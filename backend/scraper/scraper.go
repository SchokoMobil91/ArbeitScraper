package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func ScrapeJobDetail(ctx context.Context, jobLink string) (fullDesc, appLink, phone, mail string) {
	var html string
	detailCtx, detailCancel := context.WithTimeout(ctx, 30*time.Second)
	defer detailCancel()

	phone = ""
	mail = ""
	appLink = ""
	fullDesc = ""

	err := chromedp.Run(detailCtx,
		chromedp.Navigate(jobLink),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML("body", &html),
	)

	if err != nil {
		log.Printf("Error scraping detail page %s: %v", jobLink, err)
		return "", "", "", ""
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	fullDescriptionSelection := doc.Find(".ba-stelle-details__description-content")
	if fullDescriptionSelection.Length() > 0 {
		fullDesc = strings.TrimSpace(fullDescriptionSelection.Text())
	}

	appLinkSelection := doc.Find(".ba-apply-info a[target='_blank']")

	if appLinkSelection.Length() > 0 {
		link, exists := appLinkSelection.Attr("href")
		if exists && !strings.HasPrefix(link, "javascript:") {
			appLink = link
		}
	} else {
		appLink = jobLink
	}

	contactSection := doc.Find(".ba-detail-contact").Text()

	if contactSection == "" {
		contactSection = fullDesc
	}

	phone = extractTelephone(contactSection)
	mail = extractEmail(contactSection)

	return fullDesc, appLink, phone, mail
}

func ScrapeJobs() ([]Job, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var html string
	url := "https://www.arbeitsagentur.de/jobsuche/suche?angebotsart=1"

	log.Println("Scraping list page...")

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("body", &html),
	)

	if err != nil {
		return nil, fmt.Errorf("error scraping list page: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	jobs := []Job{}
	tempJobs := []Job{}

	doc.Find("li.ba-tile").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3.titel-lane span:nth-child(2)").Text())

		company := strings.TrimSpace(s.Find("h4.firma-lane").Text())
		company = strings.TrimPrefix(company, "Arbeitgeber:")
		company = strings.TrimSpace(company)

		location := strings.TrimSpace(s.Find("span.ba-icon-location-full span:nth-child(2)").Text())

		salary := strings.TrimSpace(s.Find("span.ba-icon-money span:nth-child(2)").Text())

		startDate := strings.TrimSpace(s.Find("span.ba-icon-calendar span:nth-child(2)").Text())

		shortDescription := strings.TrimSpace(s.Find("p.ba-description").Text())

		refNr := strings.TrimSpace(s.Find("span.ba-job-ref").Text())

		link, _ := s.Find("a").Attr("href")
		if strings.HasPrefix(link, "javascript:") {
			link = ""
		} else if !strings.HasPrefix(link, "http") {
			link = "https://www.arbeitsagentur.de" + link
		}

		if refNr == "" {
            refNr = link
        }

		tempJobs = append(tempJobs, Job{
			Profession:       title,
			Company:          company,
			Location:         location,
			Salary:           salary,
			StartDate:        startDate,
			ShortDescription: shortDescription,
			RefNr:            refNr,
			ExternalLink:     link,
		})
	})

	log.Printf("Scraped %d basic jobs. Starting detail scraping...", len(tempJobs))

	for i, job := range tempJobs {
		if job.ExternalLink == "" {
			continue
		}

		fullDesc, appLink, phone, mail := ScrapeJobDetail(ctx, job.ExternalLink)

		job.FullDescription = fullDesc
		job.ApplicationLink = appLink
		job.Telephone = phone
		job.Email = mail

		jobs = append(jobs, job)

		time.Sleep(1 * time.Second)

		log.Printf("Completed detail scraping for job %d/%d: %s", i+1, len(tempJobs), job.Profession)
	}

	log.Printf("Scraped %d jobs", len(jobs))
	return jobs, nil
}
