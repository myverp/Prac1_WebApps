package site

import (
	"html"
	"html/template"
	"io"
	"strings"
)

type Option struct {
	Value string
	Label string
}

type Field struct {
	Label    string
	Name     string
	Value    string
	Type     string
	Step     string
	Unit     string
	ReadOnly bool
	Options  []Option
}

type Metric struct {
	Label string
	Value string
}

type Section struct {
	Title string
	HTML  template.HTML
}

type PageData struct {
	Title         string
	Practice      string
	Breadcrumb    string
	Lead          string
	SelectedValue string
	Tags          []string
	Highlights    []string
	Fields        []Field
	ExtraForm     template.HTML
	Metrics       []Metric
	Sections      []Section
	Footer        string
}

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="uk">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>{{.Title}}</title>
  <style>
    :root {
      --bg: #f5f4ef;
      --surface: #ffffff;
      --text: #15202b;
      --muted: #607086;
      --border: rgba(21, 32, 43, 0.1);
      --accent: #244c7c;
      --shadow: 0 4px 14px rgba(21, 32, 43, 0.05);
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font: 16px/1.5 "Segoe UI", Arial, sans-serif;
    }
    a { color: inherit; }
    .page {
      width: min(1120px, calc(100% - 32px));
      margin: 0 auto;
      padding: 24px 0 56px;
    }
    .topbar {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      align-items: start;
      margin-bottom: 18px;
    }
    .brand {
      display: grid;
      gap: 3px;
    }
    .brand strong {
      font-size: 1rem;
    }
    .brand span, .breadcrumbs, .lead, .section-head p, .highlight, .note, .footer {
      color: var(--muted);
    }
    .breadcrumbs { font-size: 0.92rem; }
    .hero {
      display: grid;
      grid-template-columns: minmax(0, 1.4fr) minmax(280px, 0.8fr);
      gap: 16px;
      margin-bottom: 16px;
    }
    .card {
      background: var(--surface);
      border: 1px solid var(--border);
      border-radius: 16px;
      box-shadow: var(--shadow);
      padding: 18px;
    }
    h1, h2, h3 { margin: 0; line-height: 1.15; }
    h1 { font-size: clamp(2rem, 4vw, 3rem); max-width: 13ch; }
    h2 { font-size: 1.15rem; }
    .eyebrow {
      margin: 0 0 10px;
      color: var(--accent);
      text-transform: uppercase;
      letter-spacing: 0.12em;
      font-size: 0.78rem;
      font-weight: 700;
    }
    .lead { margin: 12px 0 0; max-width: 60ch; }
    .tags { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 14px; }
    .tag {
      display: inline-flex;
      align-items: center;
      padding: 6px 10px;
      border-radius: 999px;
      border: 1px solid var(--border);
      background: #f8f8f6;
      color: var(--muted);
      font-size: 0.88rem;
    }
    .highlights {
      margin: 0;
      padding-left: 18px;
      display: grid;
      gap: 10px;
    }
    .highlights li { margin: 0; }
    .grid {
      display: grid;
      gap: 16px;
    }
    .section-head { display: flex; justify-content: space-between; gap: 12px; align-items: start; margin-bottom: 14px; }
    .section-head p { margin: 6px 0 0; }
    .field-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 12px;
    }
    .field {
      display: grid;
      gap: 6px;
    }
    label {
      font-size: 0.9rem;
      font-weight: 600;
    }
    input, select {
      width: 100%;
      border: 1px solid rgba(21, 32, 43, 0.16);
      border-radius: 10px;
      padding: 10px 12px;
      background: #fff;
      color: var(--text);
      font: inherit;
    }
    input[readonly] { background: #f9f9f7; }
    .field small { color: var(--muted); }
    .button-row {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin-top: 14px;
    }
    .button, .button-link {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border-radius: 999px;
      padding: 10px 14px;
      text-decoration: none;
      font: inherit;
      font-weight: 700;
      border: 0;
      cursor: pointer;
    }
    .button { background: var(--accent); color: #fff; }
    .button-link { background: #eef1f5; color: var(--text); }
    .metrics {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
      gap: 12px;
    }
    .metric {
      border: 1px solid var(--border);
      border-radius: 12px;
      padding: 12px 14px;
      background: #fff;
    }
    .metric span {
      display: block;
      font-size: 0.85rem;
      color: var(--muted);
      margin-bottom: 3px;
    }
    .metric strong { font-size: 1rem; }
    .table-wrap {
      overflow-x: auto;
      border: 1px solid var(--border);
      border-radius: 12px;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 560px;
      background: #fff;
    }
    th, td {
      padding: 10px 12px;
      border-bottom: 1px solid rgba(21, 32, 43, 0.08);
      text-align: left;
      vertical-align: top;
    }
    th {
      background: #f3f6f9;
      font-size: 0.84rem;
      text-transform: uppercase;
      letter-spacing: 0.06em;
      color: var(--muted);
    }
    tr:last-child td { border-bottom: 0; }
    .section-text {
      color: var(--muted);
    }
    .section-table {
      margin-top: 10px;
    }
    .footer {
      margin-top: 16px;
      font-size: 0.9rem;
    }
    @media (max-width: 900px) {
      .hero { grid-template-columns: 1fr; }
      .topbar { flex-direction: column; }
    }
    @media (max-width: 640px) {
      .page { width: min(calc(100% - 18px), 1120px); padding-top: 16px; }
      .card { padding: 16px; }
    }
  </style>
</head>
<body>
  <main class="page">
    <div class="topbar">
      <div class="brand">
        <strong>{{.Practice}}</strong>
        <span>Програмування вебзастосунків</span>
      </div>
      <div class="breadcrumbs">{{.Breadcrumb}}</div>
    </div>

    <section class="hero">
      <div class="card">
        <p class="eyebrow">{{.Breadcrumb}}</p>
        <h1>{{.Title}}</h1>
        <p class="lead">{{.Lead}}</p>
        <div class="tags">
          {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
        </div>
      </div>
      <aside class="card">
        <div class="section-head">
          <div>
            <h2>Коротко</h2>
            <p>Ключові акценти поточної роботи.</p>
          </div>
        </div>
        <ul class="highlights">
          {{range .Highlights}}<li class="highlight">{{.}}</li>{{end}}
        </ul>
      </aside>
    </section>

    <section class="grid">
      <article class="card">
        <div class="section-head">
          <div>
            <h2>Вхідні дані</h2>
            <p>Введіть значення та натисніть обчислення.</p>
          </div>
        </div>
        <form method="get">
          <div class="field-grid">
            {{range .Fields}}
              <div class="field">
                <label for="{{.Name}}">{{.Label}}</label>
                {{if eq .Type "select"}}
                  <select id="{{.Name}}" name="{{.Name}}">
                    {{range .Options}}
                      <option value="{{.Value}}" {{if eq $.SelectedValue .Value}}selected{{end}}>{{.Label}}</option>
                    {{end}}
                  </select>
                {{else}}
                  <input id="{{.Name}}" name="{{.Name}}" type="text" inputmode="decimal" value="{{.Value}}" {{if .ReadOnly}}readonly{{end}} />
                {{end}}
                {{if .Unit}}<small>{{.Unit}}</small>{{end}}
              </div>
            {{end}}
          </div>
          {{.ExtraForm}}
          <div class="button-row">
            <button class="button" type="submit">Обчислити</button>
            <a class="button-link" href="./">Скинути</a>
          </div>
        </form>
      </article>

      <article class="card">
        <div class="section-head">
          <div>
            <h2>Результати</h2>
            <p>Поточний розрахунок для введених даних.</p>
          </div>
        </div>
        <div class="metrics">
          {{range .Metrics}}
            <div class="metric">
              <span>{{.Label}}</span>
              <strong>{{.Value}}</strong>
            </div>
          {{end}}
        </div>
      </article>

      {{range .Sections}}
      <article class="card">
        <div class="section-head">
          <div>
            <h2>{{.Title}}</h2>
          </div>
        </div>
        <div class="section-text">{{.HTML}}</div>
      </article>
      {{end}}
    </section>

    <div class="footer">{{.Footer}}</div>
  </main>
</body>
</html>`))

func Render(w io.Writer, data PageData) error {
	return pageTemplate.Execute(w, data)
}

func Table(headers []string, rows [][]string) template.HTML {
	var b strings.Builder
	b.WriteString(`<div class="table-wrap section-table"><table><thead><tr>`)
	for _, h := range headers {
		b.WriteString("<th>")
		b.WriteString(html.EscapeString(h))
		b.WriteString("</th>")
	}
	b.WriteString("</tr></thead><tbody>")
	for _, row := range rows {
		b.WriteString("<tr>")
		for _, cell := range row {
			b.WriteString("<td>")
			b.WriteString(html.EscapeString(cell))
			b.WriteString("</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table></div>")
	return template.HTML(b.String())
}
