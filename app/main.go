package main

import (
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"log"
	"github.com/pamungkaski/go-fuzzy-logic"
	"strconv"
	"strings"
	"sort"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"gonum.org/v1/plot/vg"
	"fmt"
	"gonum.org/v1/plot/vg/draw"
	"math"
)

func main() {
	///// DATA READING
	csvFile, _ := os.Open("DataTugas2.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	var data []fuzzy.FuzzyNumber
	reader.Read()
	for  {
		var fn fuzzy.FuzzyNumber
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fn.Family.Number = line[0]
		fn.Family.Income, err = strconv.ParseFloat(strings.Replace(line[1], " ", "", -1), 64)
		if err != nil {
			log.Fatal(err)
		}
		fn.Family.Debt, err = strconv.ParseFloat(strings.Replace(line[2], " ", "", -1), 64)
		if err != nil {
			log.Fatal(err)
		}

		fn.Family.Income *= 1000
		fn.Family.Debt *= 1000
		data = append(data, fn)
	}

	file, _ := os.Create("TebakanTugas2.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	////// FUZI
	blt := fuzzy.BLT{}
	for i := range data {
		blt.Fuzzification(&data[i])
		blt.Inference(&data[i])
		blt.Defuzzification(&data[i])
	}

	///// SORT DATA
	sort.Slice(data, func(i, j int) bool {
		return data[i].CrispValue > data[j].CrispValue
	})

	///////// INSERTING DATA
	head := []string{
		"No",
		"Income",
		"Debt",
		"Crisp",
	}
	if err := writer.Write(head); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	for i := range data[:20] {
		csvData := []string{
			fmt.Sprintf("%s", data[i].Family.Number),
			fmt.Sprintf("%.0f", data[i].Family.Income),
			fmt.Sprintf("%.0f", data[i].Family.Debt),
			fmt.Sprintf("%f", data[i].CrispValue),
		}
		if err := writer.Write(csvData); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	///////// PLOT DATA AND RESULT

	scatterData := make(plotter.XYZs, len(data))
	for i := range data {
		scatterData[i].X = data[i].Family.Income
		scatterData[i].Y = data[i].Family.Debt
		scatterData[i].Z = data[i].CrispValue
	}
	sort.Slice(scatterData, func(i, j int) bool {
		return scatterData[i].Z < scatterData[j].Z
	})

	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, xyz := range scatterData {
		if xyz.Z > maxZ {
			maxZ = xyz.Z
		}
		if xyz.Z < minZ {
			minZ = xyz.Z
		}
	}

	plt, _ := plot.New()
	plt.Title.Text = "Fuzzy Logic BLT Forecast"
	plt.X.Label.Text = "Income"
	plt.Y.Label.Text = "Debt"
	plt.Legend.Add("The bigger the circle, the bigger the value")

	sc, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	// Specify style for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		x, y, z := scatterData.XYZ(i)
		c := color.RGBA{R: 190 + uint8(x*20), G: 128 + uint8(y*20), B: 0 + uint8(z*20), A: 255}
		var minRadius, maxRadius = vg.Points(-10), vg.Points(10)
		rng := maxRadius - minRadius
		d := (z - minZ) / (maxZ - minZ)
		r := vg.Length(d) * rng * 2
		return draw.GlyphStyle{Color: c, Radius: r, Shape: draw.CircleGlyph{}}
	}
	plt.Add(sc)

	plt.Save(15*vg.Inch, 15*vg.Inch, "contour.svg")


	/////// PLOT FUNCTION INCOME

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Income Memberships Functions"
	p.X.Label.Text = "Income"
	p.Y.Label.Text = "Fuzz"

	low := plotter.NewFunction(blt.IncomeLow)
	low.Color = color.RGBA{R: 255, A: 255}

	mid := plotter.NewFunction(blt.IncomeMiddle)
	mid.Color = color.RGBA{G: 255, A: 255}

	hig := plotter.NewFunction(blt.IncomeHigh)
	hig.Color = color.RGBA{B: 255, A: 255}

	p.Add(low, mid, hig)

	// Set the axis ranges.  Unlike other data sets,
	// functions don't set the axis ranges automatically
	// since functions don't necessarily have a
	// finite range of x and y values.
	p.X.Min = 0
	p.X.Max = 2000
	p.Y.Min = 0
	p.Y.Max = 1

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "functionsIncome.svg"); err != nil {
		panic(err)
	}

	///////// PLOT FUNCTION DEBT
	p, err = plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Debt Membership Functions"
	p.X.Label.Text = "Debt"
	p.Y.Label.Text = "Fuzz"

	low = plotter.NewFunction(blt.DebtLow)
	low.Color = color.RGBA{R: 255, A: 255}

	mid = plotter.NewFunction(blt.DebtMiddle)
	mid.Color = color.RGBA{G: 255, A: 255}

	hig = plotter.NewFunction(blt.DebtHigh)
	hig.Color = color.RGBA{B: 255, A: 255}

	p.Add(low, mid, hig)

	// Set the axis ranges.  Unlike other data sets,
	// functions don't set the axis ranges automatically
	// since functions don't necessarily have a
	// finite range of x and y values.
	p.X.Min = 0
	p.X.Max = 100000
	p.Y.Min = 0
	p.Y.Max = 1

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "functionsDebt.svg"); err != nil {
		panic(err)
	}
}