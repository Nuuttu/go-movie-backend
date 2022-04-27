package excelporter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	mystructs "example.com/mysctructs"
	"github.com/xuri/excelize/v2"
)

func Main() {
	fmt.Println("Excel importer")
}

type User = mystructs.User
type Watch = mystructs.Watch
type Movie = mystructs.Movie

func Excelimporter() {

	f, err := excelize.OpenFile("Medialists.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	/*
		cell, err := f.GetCellValue("Movies", "A4")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cell)
	*/
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Movies")
	if err != nil {
		fmt.Println(err)
		return
	}
	for is, row := range rows {
		if is > 3 {
			//					for ir, colCell := range row {
			// }
			var newMovie Movie
			newMovie.SetId()
			newMovie.Name = row[0]
			newMovie.Year = 0
			rating, e := strconv.ParseFloat(row[1], 10)
			if e != nil {
				fmt.Println("e", e)
			}
			newMovie.Rating = int(rating)
			if len(row) == 5 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n %s,\n %s,\n", row[0], row[1], row[2], row[3], row[4])

				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				if len(row[3]) > 0 {
					//fmt.Printf("\n\n-%s-\n", row[3])
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[3], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
				if len(row[4]) > 0 {
					//fmt.Printf("\n\n-%s-\n", strings.Replace(row[4], ".", "-", -1))
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[4], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			} else if len(row) == 4 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n %s,\n", row[0], row[1], row[2], row[3])
				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				if len(row[3]) > 0 {
					//fmt.Println(row[3])
					var newWatch Watch
					wDate, _ := time.Parse("2006-01-02", strings.Replace(row[3], ".", "-", -1))
					//fmt.Println("wDate", wDate)
					newWatch.Date = wDate
					newWatch.SetId()
					//fmt.Println("enwWathc", newWatch)
					newMovie.Watches = append(newMovie.Watches, newWatch)
				}
			} else if len(row) == 3 {
				//fmt.Printf("is?\n %s,\n %s,\n %s,\n", row[0], row[1], row[2])
				if len(row[2]) > 0 {
					//fmt.Println(row[2])
					newMovie.Review = row[2]
				}
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			} else {
				//fmt.Printf("is?\n %s,\n %s,\n", row[0], row[1])
				var newWatch Watch
				newWatch.SetId()
				newMovie.Watches = append(newMovie.Watches, newWatch)
			}

			//fmt.Println("newMovie", newMovie)
			//movieList = append(movieList, newMovie)
			//fmt.Println()
		}
	}

}
