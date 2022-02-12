package pdf

import (
	"fmt"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/util"
)

type PDFConstructor struct {
	Invoice    *db.Invoice
	OutputPath string
}

func Construct(p PDFConstructor) error {

	blueColor := getBlueColor()

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)

	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				m.Text(p.Invoice.User.GetFullName(), props.Text{
					Size:        12,
					Align:       consts.Left,
					Style:       consts.Bold,
					Extrapolate: false,
				})
				m.Text(fmt.Sprintf("Tel: %s", p.Invoice.User.Phone), props.Text{
					Top:   12,
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
		getBillTo(m, p.Invoice.Client)
		m.ColSpace(6)
		getInvoiceDetails(m, p.Invoice)
	})

	getTable(m, p.Invoice)

	m.Row(5, func() {
		m.Col(9, func() {
			m.Text("Subtotal (exc. GST):", props.Text{
				Top:   5,
				Size:  10,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text(
				fmt.Sprintf("$%.2f", p.Invoice.GetSubtotal()), props.Text{
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
			m.Text(
				fmt.Sprintf("$%.2f", p.Invoice.GetGST()),
				props.Text{
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
			m.Text(
				fmt.Sprintf("$%.2f", p.Invoice.GetTotal()),
				props.Text{
					Top:   5,
					Size:  12,
					Align: consts.Right,
				})
		})
	})

	if err := m.OutputFileAndClose(p.OutputPath); err != nil {
		return errors.Annotate(err, "could not save file")
	}

	return nil
}

func getBillTo(m pdf.Maroto, client *db.User) {
	m.Col(3, func() {
		m.Text("Bill to", props.Text{
			Top:         3,
			Size:        8,
			Align:       consts.Left,
			Style:       consts.Bold,
			Extrapolate: false,
		})
		m.Text(
			client.GetFullName(),
			props.Text{
				Top:         6,
				Size:        8,
				Align:       consts.Left,
				Style:       consts.Italic,
				Extrapolate: false,
			})
		m.Text(client.Company, props.Text{
			Top:         9,
			Size:        8,
			Align:       consts.Left,
			Extrapolate: false,
		})
		m.Text(
			fmt.Sprintf("%s", client.Address.FirstLine),
			props.Text{
				Top:         12,
				Size:        8,
				Align:       consts.Left,
				Extrapolate: false,
			})
		m.Text(
			fmt.Sprintf("%s", client.Address.SecondLine),
			props.Text{
				Top:         15,
				Size:        8,
				Align:       consts.Left,
				Extrapolate: false,
			})
		m.Text(
			fmt.Sprintf("%s", client.Address.Postcode),
			props.Text{
				Top:         18,
				Size:        8,
				Align:       consts.Left,
				Extrapolate: false,
			})
		m.Text(
			fmt.Sprintf("%s", client.Address.Country),
			props.Text{
				Top:         21,
				Size:        8,
				Align:       consts.Left,
				Extrapolate: false,
			})
	})
}

func getInvoiceDetails(m pdf.Maroto, i *db.Invoice) {
	m.Col(3, func() {
		m.Text("Invoice number", props.Text{
			Top:         3,
			Size:        8,
			Align:       consts.Right,
			Style:       consts.Bold,
			Extrapolate: false,
		})
		m.Text(
			fmt.Sprintf("%d", i.Number),
			props.Text{
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
		m.Text(
			util.ToFormattedDate(i.IssueDate),
			props.Text{
				Top:         12,
				Size:        8,
				Align:       consts.Right,
				Extrapolate: false,
			})
		m.Text("Due Date", props.Text{
			Top:         15,
			Size:        8,
			Align:       consts.Right,
			Extrapolate: false,
		})
		m.Text(
			util.ToFormattedDate(i.DueDate),
			props.Text{
				Top:         18,
				Size:        8,
				Align:       consts.Right,
				Extrapolate: false,
			})
	})
}

func getTable(m pdf.Maroto, i *db.Invoice) {

	grayColor := getGrayColor()

	m.SetBackgroundColor(getDarkGrayColor())
	m.Row(5, func() {
		m.ColSpace(12)
	})
	m.SetBackgroundColor(color.NewWhite())

	m.TableList(getHeader(), getContents(i), props.TableList{
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

func getContents(i *db.Invoice) [][]string {
	return [][]string{
		{i.Description, fmt.Sprintf("%.2f", i.Hours), fmt.Sprintf("%.2f", i.Rate)},
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
