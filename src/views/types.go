package views

import "github.com/rivo/tview"

type Element struct {
	View       tview.Primitive
	ParentFlex *tview.Flex
}

const (
	// Class Detail View box names
	ClassDetailsGeneralStatsBox = "generalStats"
	ClassDetailsEpicStatsBox    = "epicStats"
	ClassDetailsFooterBox       = "footer"
	ClassDetailsPrereqBox       = "prereq"
	ClassDetailsSpellsBox       = "spells"
)
