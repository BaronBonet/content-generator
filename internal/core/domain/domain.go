package domain

import "time"

type NewsArticle struct {
	Title  string
	Body   string
	Date   Date
	Url    string
	Source string
}

type Date struct {
	Day   int
	Month time.Month
	Year  int
}

type ImagePath string
