/*

Copyright (c) 2018 sec.xiaomi.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package routers

import (
	"gopkg.in/macaron.v1"
	"x-patrol/models"
	"x-patrol/vars"

	"github.com/go-macaron/session"
	"net/url"
	"strconv"
	"x-patrol/util/common"
)

func GetDetailedReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		codeResultDetail, _ := models.GetCodeResultDetailById(int64(Id))
		ctx.Data["detailed_report"] = codeResultDetail
		ctx.HTML(200, "report_detail")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ListGithubSearchResultByStatus(ctx *macaron.Context, sess session.Store) {
	status, _ := strconv.Atoi(ctx.Params(":status"))
	ctx.SetCookie("status", strconv.Itoa(status))
	p := 0

	renderDataForGithubSearchResult(ctx, sess, p, status)
}

func ListGithubSearchResult(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	status := 0

	if ctx.GetCookie("status") != "" {
		status, _ = strconv.Atoi(ctx.GetCookie("status"))
	}
	renderDataForGithubSearchResult(ctx, sess, p, status)
}

func renderDataForGithubSearchResult(ctx *macaron.Context, sess session.Store, p, status int) {
	if sess.Get("admin") != nil {
		p, pre, next := common.GetPreAndNext(p)
		reports, pages, count := models.ListGithubSearchResultPage(p, status)
		pageList := common.GetPageList(p, vars.PageStep, pages)
		lastPage := 0
		if len(pageList) >= 1 {
			lastPage = pageList[len(pageList)-1]
		}

		ctx.Data["reports"] = reports
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["status"] = status
		ctx.Data["count"] = count
		ctx.Data["lastPage"] = lastPage
		ctx.HTML(200, "report_github")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func getRefer(ctx *macaron.Context) string {
	refer := "/admin/reports/github/"
	if _, ok := ctx.Req.Header["Referer"]; len(ctx.Req.Header["Referer"]) > 0 && ok {
		u := ctx.Req.Header["Referer"][0]
		urlParsed, err := url.Parse(u)
		if err == nil {
			refer = urlParsed.RequestURI()
		}
	}
	return refer
}

func ConfirmReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ConfirmReportById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func CancelReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.CancelReportById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepoById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		omitRepo := true
		has, result, err := models.GetReportById(int64(Id), omitRepo)
		if err == nil && has {
			models.DisableRepoByUrl(result.Repository.GetHTMLURL())
			models.CancelReportsByRepo(int64(Id))
		}
		models.CancelReportById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}

/*
For local code search
*/

func ListLocalSearchResult(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		reports, _ := models.ListSearchResult()
		ctx.Data["reports"] = reports
		ctx.HTML(200, "report_search")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ListLocalSearchResultPage(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pre := p - 1
	if pre <= 0 {
		pre = 1
	}
	next := p + 1

	if sess.Get("admin") != nil {
		reports, pages, _ := models.ListSearchResultPage(p)
		pList := 0
		if pages-p > 10 {
			pList = p + 10
		} else {
			pList = pages
		}

		pageList := make([]int, 0)
		if pages <= 10 {
			for i := 1; i <= pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			if p <= 10 {
				for i := 1; i <= pList; i++ {
					pageList = append(pageList, i)
				}
			} else {
				t := p + 5
				if t > pages {
					t = pages
				}
				for i := p - 5; i <= t; i++ {
					pageList = append(pageList, i)
				}
			}
		}

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["reports"] = reports
		ctx.HTML(200, "report_search")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ConfirmSearchResultById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ConfirmSearchResultById(int64(Id))
		ctx.Redirect(getRefer(ctx))
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func CancelSearchResultById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.CancelSearchResultById(int64(Id))
		ctx.Redirect("/admin/reports/search/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableSearchRepoById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		has, result, err := models.GetSearchResultById(int64(Id))
		if err == nil && has {
			repoUrl := result.Repo
			models.DisableRepoByUrl(repoUrl)
		}
		models.CancelSearchResultById(int64(Id))
		ctx.Redirect("/admin/reports/search/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
