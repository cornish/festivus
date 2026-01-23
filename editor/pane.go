package editor

// Pane represents a single view into a document.
// Each pane maintains its own scroll position for independent viewing.
type Pane struct {
	documentIdx int // Index into Editor.documents
	scrollY     int // Pane-specific vertical scroll position
	scrollX     int // Pane-specific horizontal scroll position
}

// NewPane creates a new pane viewing the document at the given index.
func NewPane(docIdx int) *Pane {
	return &Pane{
		documentIdx: docIdx,
		scrollY:     0,
		scrollX:     0,
	}
}

// DocumentIdx returns the index of the document this pane is viewing.
func (p *Pane) DocumentIdx() int {
	return p.documentIdx
}

// SetDocumentIdx sets the document index for this pane.
func (p *Pane) SetDocumentIdx(idx int) {
	p.documentIdx = idx
}

// ScrollY returns the vertical scroll position.
func (p *Pane) ScrollY() int {
	return p.scrollY
}

// SetScrollY sets the vertical scroll position.
func (p *Pane) SetScrollY(y int) {
	if y < 0 {
		y = 0
	}
	p.scrollY = y
}

// ScrollX returns the horizontal scroll position.
func (p *Pane) ScrollX() int {
	return p.scrollX
}

// SetScrollX sets the horizontal scroll position.
func (p *Pane) SetScrollX(x int) {
	if x < 0 {
		x = 0
	}
	p.scrollX = x
}
