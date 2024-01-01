package accum

// LessonCounter is a visitor that merely counts lessons.
type LessonCounter struct {
	count int
}

// NewTutorialLessonCounter makes a new LessonCounter.
func NewTutorialLessonCounter() *LessonCounter {
	return &LessonCounter{0}
}

// Count is the reason this visitor exists.
func (v *LessonCounter) Count() int {
	return v.count
}

// VisitCodeBlock does nothing.
func (v *LessonCounter) VisitCodeBlock(b *CodeBlock) {
}

// VisitLessonDoc increments the count.
func (v *LessonCounter) VisitLessonDoc(l *LessonDoc) {
	v.count++
}

// VisitCourse visits children.
func (v *LessonCounter) VisitCourse(c *Course) {
	for _, x := range c.Children() {
		x.Accept(v)
	}
}
