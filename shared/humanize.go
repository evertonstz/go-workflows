package shared

import (
	"fmt"
	"time"

	"github.com/evertonstz/go-workflows/shared/di"
	"github.com/evertonstz/go-workflows/shared/di/services"
)

type I18nTranslator interface {
	Translate(key string) string
}

func HumanizeTime(t time.Time) string {
	i18n := di.GetService[*services.I18nService](di.I18nServiceKey)
	return HumanizeTimeWithService(t, i18n)
}

func HumanizeTimeWithService(t time.Time, i18n I18nTranslator) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 0 {
		diff = -diff
		return humanizeTimeFutureWithService(diff, i18n)
	}

	return humanizeTimePastWithService(diff, i18n)
}

func humanizeTimePastWithService(diff time.Duration, i18n I18nTranslator) string {
	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := days / 30
	years := days / 365

	switch {
	case seconds < 60:
		if seconds <= 1 {
			return i18n.Translate("time_just_now")
		}
		return fmt.Sprintf(i18n.Translate("time_seconds_ago"), seconds)
	case minutes < 60:
		if minutes == 1 {
			return i18n.Translate("time_minute_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_minutes_ago"), minutes)
	case hours < 24:
		if hours == 1 {
			return i18n.Translate("time_hour_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_hours_ago"), hours)
	case days < 7:
		if days == 1 {
			return i18n.Translate("time_day_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_days_ago"), days)
	case weeks < 4:
		if weeks == 1 {
			return i18n.Translate("time_week_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_weeks_ago"), weeks)
	case months < 12:
		if months == 1 {
			return i18n.Translate("time_month_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_months_ago"), months)
	default:
		if years == 1 {
			return i18n.Translate("time_year_ago")
		}
		return fmt.Sprintf(i18n.Translate("time_years_ago"), years)
	}
}

func humanizeTimeFutureWithService(diff time.Duration, i18n I18nTranslator) string {
	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	switch {
	case seconds < 60:
		return i18n.Translate("time_in_seconds")
	case minutes < 60:
		if minutes == 1 {
			return i18n.Translate("time_in_minute")
		}
		return fmt.Sprintf(i18n.Translate("time_in_minutes"), minutes)
	case hours < 24:
		if hours == 1 {
			return i18n.Translate("time_in_hour")
		}
		return fmt.Sprintf(i18n.Translate("time_in_hours"), hours)
	default:
		if days == 1 {
			return i18n.Translate("time_in_day")
		}
		return fmt.Sprintf(i18n.Translate("time_in_days"), days)
	}
}
