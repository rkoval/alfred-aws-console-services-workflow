package searchers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	aw "github.com/deanishe/awgo"
)

type Searcher interface {
	Search(wf *aw.Workflow, query string, session *session.Session, forceFetch bool, fullQuery string) error
}
