package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"webapps/internal/site"
)

type fuelPreset struct {
	Label  string
	Values map[string]float64
}

var presets = map[string]fuelPreset{
	"control": {
		Label:  "Контрольний приклад",
		Values: map[string]float64{"h": 1.9, "c": 21.1, "s": 2.6, "n": 0.2, "o": 7.1, "w": 53.0, "a": 14.1},
	},
	"0": {Label: "Варіант 0", Values: map[string]float64{"h": 3.2, "c": 54.4, "s": 2.3, "n": 1.0, "o": 3.1, "w": 20.0, "a": 16.0}},
	"1": {Label: "Варіант 1", Values: map[string]float64{"h": 3.7, "c": 50.6, "s": 4.0, "n": 1.1, "o": 8.0, "w": 13.0, "a": 19.6}},
	"2": {Label: "Варіант 2", Values: map[string]float64{"h": 4.2, "c": 62.1, "s": 3.0, "n": 1.2, "o": 0.6, "w": 40.7, "a": 15.8}},
	"3": {Label: "Варіант 3", Values: map[string]float64{"h": 3.8, "c": 62.4, "s": 3.6, "n": 1.1, "o": 4.3, "w": 6.0, "a": 18.8}},
	"4": {Label: "Варіант 4", Values: map[string]float64{"h": 3.4, "c": 70.6, "s": 2.7, "n": 1.2, "o": 1.9, "w": 5.0, "a": 15.2}},
	"5": {Label: "Варіант 5", Values: map[string]float64{"h": 2.8, "c": 72.3, "s": 2.0, "n": 1.1, "o": 1.3, "w": 5.5, "a": 15.0}},
	"6": {Label: "Варіант 6", Values: map[string]float64{"h": 1.5, "c": 76.4, "s": 1.7, "n": 0.8, "o": 1.3, "w": 5.0, "a": 13.3}},
	"7": {Label: "Варіант 7", Values: map[string]float64{"h": 1.4, "c": 71.7, "s": 1.8, "n": 0.8, "o": 1.4, "w": 6.0, "a": 16.9}},
	"8": {Label: "Варіант 8", Values: map[string]float64{"h": 1.4, "c": 70.5, "s": 1.7, "n": 0.8, "o": 1.9, "w": 7.0, "a": 16.7}},
	"9": {Label: "Варіант 9", Values: map[string]float64{"h": 2.6, "c": 38.6, "s": 3.8, "n": 0.8, "o": 3.1, "w": 11.0, "a": 40.1}},
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Prac1 running on http://127.0.0.1:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	presetKey := r.URL.Query().Get("preset")
	if presetKey == "" {
		presetKey = "control"
	}
	preset, ok := presets[presetKey]
	if !ok {
		presetKey = "control"
		preset = presets[presetKey]
	}

	values := make(map[string]float64, len(preset.Values))
	for key, value := range preset.Values {
		values[key] = value
	}
	for _, name := range []string{"h", "c", "s", "n", "o", "w", "a"} {
		if raw := r.URL.Query().Get(name); raw != "" {
			values[name] = parseFloat(raw, values[name])
		}
	}

	h := values["h"]
	c := values["c"]
	s := values["s"]
	n := values["n"]
	o := values["o"]
	water := values["w"]
	ash := values["a"]

	baseDry := 100 - water
	baseComb := 100 - water - ash
	kRS := 0.0
	if baseDry > 0 {
		kRS = 100 / baseDry
	}
	kRG := 0.0
	if baseComb > 0 {
		kRG = 100 / baseComb
	}

	qri := 0.339*c + 1.03*h - 0.1088*(o-s) - 0.025*water
	qd := 0.0
	if baseDry > 0 {
		qd = qri * 100 / baseDry
	}
	qdaf := 0.0
	if baseComb > 0 {
		qdaf = qri * 100 / baseComb
	}

	dryRows := [][]string{
		{"Hc", fmtPercent(h * kRS)},
		{"Cc", fmtPercent(c * kRS)},
		{"Sc", fmtPercent(s * kRS)},
		{"Nc", fmtPercent(n * kRS)},
		{"Oc", fmtPercent(o * kRS)},
		{"Ac", fmtPercent(ash * kRS)},
		{"Разом", fmtPercent((h + c + s + n + o + ash) * kRS)},
	}
	combRows := [][]string{
		{"Hg", fmtPercent(h * kRG)},
		{"Cg", fmtPercent(c * kRG)},
		{"Sg", fmtPercent(s * kRG)},
		{"Ng", fmtPercent(n * kRG)},
		{"Og", fmtPercent(o * kRG)},
		{"Разом", fmtPercent((h + c + s + n + o) * kRG)},
	}

	data := site.PageData{
		Title:         "Практична робота 1",
		Practice:      "Практична робота 1",
		Breadcrumb:    "Завдання 1",
		Lead:          "Калькулятор переводить склад робочої маси в суху та горючу і рахує нижчу теплоту згоряння за формулою Мендєлєєва.",
		SelectedValue: presetKey,
		Tags:          []string{"Перерахунок маси", "Формула Мендєлєєва", "Контрольний приклад"},
		Highlights: []string{
			"Показники робочої, сухої та горючої маси виводяться в одній формі.",
			"Контрольний приклад з PDF відтворюється без додаткових налаштувань.",
			"Поля можна редагувати вручну, а сторінка оновлюється після обчислення.",
		},
		Fields: []site.Field{
			{
				Label:   "Набір",
				Name:    "preset",
				Type:    "select",
				Options: presetOptions(),
				Value:   presetKey,
			},
			numericField("H, %", "h", h, "0.01"),
			numericField("C, %", "c", c, "0.01"),
			numericField("S, %", "s", s, "0.01"),
			numericField("N, %", "n", n, "0.01"),
			numericField("O, %", "o", o, "0.01"),
			numericField("W, %", "w", water, "0.01"),
			numericField("A, %", "a", ash, "0.01"),
		},
		Metrics: []site.Metric{
			{Label: "kRS", Value: fmtNumber(kRS, 2)},
			{Label: "kRG", Value: fmtNumber(kRG, 2)},
			{Label: "Qri", Value: fmtNumber(qri, 4)},
			{Label: "Qd", Value: fmtNumber(qd, 4)},
			{Label: "Qdaf", Value: fmtNumber(qdaf, 4)},
			{Label: "Перевірка сум", Value: fmt.Sprintf("%s / %s", fmtNumber((h+c+s+n+o+ash)*kRS, 2), fmtNumber((h+c+s+n+o)*kRG, 2))},
		},
		Sections: []site.Section{
			{
				Title: "Розклад складу",
				HTML:  site.Table([]string{"Суха маса", "Значення"}, dryRows),
			},
			{
				Title: "Горюча маса",
				HTML:  site.Table([]string{"Горюча маса", "Значення"}, combRows),
			},
			{
				Title: "Примітка до формул",
				HTML:  template.HTML(`<p>Перерахунок виконано за коефіцієнтами таблиці 1.1. Нижча теплота згоряння обчислюється у МДж/кг для робочої, сухої та горючої маси.</p>`),
			},
		},
		Footer: "Висновок: сторінка коректно перераховує склад палива та відтворює контрольний приклад із завдання.",
	}

	if err := site.Render(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func presetOptions() []site.Option {
	order := []string{"control", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	opts := make([]site.Option, 0, len(order))
	for _, key := range order {
		opts = append(opts, site.Option{Value: key, Label: presets[key].Label})
	}
	return opts
}

func numericField(label, name string, value float64, step string) site.Field {
	return site.Field{Label: label, Name: name, Value: fmtNumber(value, 2), Type: "number", Step: step}
}

func parseFloat(raw string, fallback float64) float64 {
	if v, err := strconv.ParseFloat(strings.ReplaceAll(raw, ",", "."), 64); err == nil {
		return v
	}
	return fallback
}

func fmtNumber(value float64, digits int) string {
	return strings.ReplaceAll(fmt.Sprintf("%.*f", digits, value), ".", ",")
}

func fmtPercent(value float64) string {
	return fmtNumber(value, 2) + "%"
}
