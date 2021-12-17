package pdf

import (
	"fmt"
	"os"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/rstorr/wham-platform/db"
)

type PDFConstructor struct {
	invoice db.Invoice
}

func Construct() {

	blueColor := getBlueColor()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				m.Text("Reuben Storr", props.Text{
					Size:        12,
					Align:       consts.Left,
					Style:       consts.Bold,
					Extrapolate: false,
				})
				m.Text("Tel: 027 645 8704", props.Text{
					Top:   12,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
				m.Text("www.rstorr.io", props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
					Color: blueColor,
				})
			})

			m.ColSpace(6)

		})
	})

	m.Row(30, func() {
		getBillTo(m)
		m.ColSpace(6)
		getInvoiceDetails(m)
	})

	getTable(m)

	m.Row(5, func() {
		m.Col(9, func() {
			m.Text("Subtotal (exc. GST):", props.Text{
				Top:   5,
				Size:  10,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("$9,600", props.Text{
				Top:   5,
				Size:  10,
				Align: consts.Right,
			})
		})
	})
	m.Row(5, func() {
		m.Col(9, func() {
			m.Text("GST:", props.Text{
				Top:   5,
				Size:  10,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("$1,440", props.Text{
				Top:   5,
				Size:  10,
				Align: consts.Right,
			})
		})
	})
	m.Row(5, func() {
		m.Col(9, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  12,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("$11,040.00", props.Text{
				Top:   5,
				Size:  12,
				Align: consts.Right,
			})
		})
	})

	err := m.OutputFileAndClose("assets/billing.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}
}

func getBillTo(m pdf.Maroto) {
	m.Col(3, func() {
		m.Text("Bill to", props.Text{
			Top:         3,
			Size:        8,
			Align:       consts.Left,
			Style:       consts.Bold,
			Extrapolate: false,
		})
		m.Text("Paul Swettenham", props.Text{
			Top:         6,
			Size:        8,
			Align:       consts.Left,
			Style:       consts.Italic,
			Extrapolate: false,
		})
		m.Text("Sunstone Talent", props.Text{
			Top:         9,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
		m.Text("35 Sir William Pickering Dr ", props.Text{
			Top:         12,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
		m.Text("Burnside, Christchurch", props.Text{
			Top:         15,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
		m.Text("8053", props.Text{
			Top:         18,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
		m.Text("New Zealand", props.Text{
			Top:         21,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
	})
}

func getInvoiceDetails(m pdf.Maroto) {
	m.Col(3, func() {
		m.Text("Invoice number", props.Text{
			Top:         3,
			Size:        8,
			Align:       consts.Right,
			Style:       consts.Bold,
			Extrapolate: false,
		})
		m.Text("000000009", props.Text{
			Top:         6,
			Size:        8,
			Align:       consts.Right,
			Style:       consts.Italic,
			Extrapolate: false,
		})
		m.Text("Issue Date", props.Text{
			Top:         9,
			Size:        8,
			Align:       consts.Right,
			Extrapolate: false,
		})
		m.Text("28/10/2021", props.Text{
			Top:         12,
			Size:        8,
			Align:       consts.Right,
			Extrapolate: false,
		})
		m.Text("Issue Date", props.Text{
			Top:         15,
			Size:        8,
			Align:       consts.Right,
			Extrapolate: false,
		})
		m.Text("28/10/2021", props.Text{
			Top:         18,
			Size:        8,
			Align:       consts.Right,
			Extrapolate: false,
		})
	})
}

func getTable(m pdf.Maroto) {

	grayColor := getGrayColor()

	m.SetBackgroundColor(getDarkGrayColor())
	m.Row(5, func() {
		m.ColSpace(12)
	})
	m.SetBackgroundColor(color.NewWhite())

	m.TableList(getHeader(), getContents(), props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{7, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{7, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

}

func getHeader() []string {
	return []string{"Description", "Quantity", "Amount($) ex GST"}
}

func getContents() [][]string {
	return [][]string{
		{"Billable hours 1/10/21 - 31/10/21", "120", "9600"},
	}
}

func getDarkGrayColor() color.Color {
	return color.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() color.Color {
	return color.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}
