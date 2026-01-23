package editor

// SplitOrientation defines how the editor is split.
type SplitOrientation int

const (
	SplitNone       SplitOrientation = iota // No split (single pane)
	SplitHorizontal                         // Top/bottom split
	SplitVertical                           // Left/right split
)

// SplitLayout manages the split view state.
type SplitLayout struct {
	orientation SplitOrientation
	pane1       *Pane // Top or Left pane
	pane2       *Pane // Bottom or Right pane
	activePane  int   // 0 for pane1, 1 for pane2
}

// NewSplitLayout creates a new split layout with the given orientation.
func NewSplitLayout(orientation SplitOrientation, doc1Idx, doc2Idx int) *SplitLayout {
	return &SplitLayout{
		orientation: orientation,
		pane1:       NewPane(doc1Idx),
		pane2:       NewPane(doc2Idx),
		activePane:  0,
	}
}

// Orientation returns the split orientation.
func (s *SplitLayout) Orientation() SplitOrientation {
	return s.orientation
}

// ActivePaneIndex returns the index of the active pane (0 or 1).
func (s *SplitLayout) ActivePaneIndex() int {
	return s.activePane
}

// ActivePane returns the currently active pane.
func (s *SplitLayout) ActivePane() *Pane {
	if s.activePane == 0 {
		return s.pane1
	}
	return s.pane2
}

// InactivePane returns the currently inactive pane.
func (s *SplitLayout) InactivePane() *Pane {
	if s.activePane == 0 {
		return s.pane2
	}
	return s.pane1
}

// Pane1 returns the first pane (top or left).
func (s *SplitLayout) Pane1() *Pane {
	return s.pane1
}

// Pane2 returns the second pane (bottom or right).
func (s *SplitLayout) Pane2() *Pane {
	return s.pane2
}

// SwitchPane toggles the active pane.
func (s *SplitLayout) SwitchPane() {
	if s.activePane == 0 {
		s.activePane = 1
	} else {
		s.activePane = 0
	}
}

// SetActivePane sets which pane is active (0 or 1).
func (s *SplitLayout) SetActivePane(idx int) {
	if idx == 0 || idx == 1 {
		s.activePane = idx
	}
}

// Panes returns both panes as a slice for iteration.
func (s *SplitLayout) Panes() []*Pane {
	return []*Pane{s.pane1, s.pane2}
}
