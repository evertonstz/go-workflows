package shared

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/evertonstz/go-workflows/models"
)

// ConvertItemsToBubbleList converts a slice of models.Item to a slice of list.Item
// suitable for use with bubbletea/list.Model.
func ConvertItemsToBubbleList(modelItems []models.Item) []list.Item {
	bubbleItems := make([]list.Item, len(modelItems))
	for i, mi := range modelItems {
		bubbleItems[i] = mi // models.Item now implements list.Item
	}
	return bubbleItems
}
