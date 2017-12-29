package toolkit
import (
"github.com/emirpasic/gods/maps/treemap"
	"math"
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"github.com/emirpasic/gods/utils"
)

const MISSING float64  = math.MaxFloat64;
type Matrix struct {
	m_data [][]float64
	m_attr_name []string
	m_str_to_enum []*treemap.Map
	m_enum_to_str []*treemap.Map
}

// Creates a 0x0 matrix. You should call loadARFF or setSize next.
func NewEmptyMatrix()Matrix{
	return Matrix{};
}

// Copies the specified portion of that matrix into this matrix
func NewMatrix(that Matrix, rowStart, colStart, rowCount, colCount int)Matrix {
	m := Matrix{};
	for j := 0; j < rowCount; j++ {
		rowSrc := that.Row(rowStart + j);
		rowDest := make([]float64,colCount);
		for i := 0; i < colCount; i++{
			rowDest[i] = rowSrc[colStart + i];
		}
		m.m_data = append(m.m_data,rowDest);
	}
	for i := 0; i < colCount; i++ {
		m.m_attr_name = append(m.m_attr_name,that.AttrName(colStart + i));
		m.m_str_to_enum = append(m.m_str_to_enum,that.m_str_to_enum[colStart + i]);
		m.m_enum_to_str = append(m.m_enum_to_str,that.m_enum_to_str[colStart + i]);
	}
	return m;
}

// Adds a copy of the specified portion of that matrix to this matrix
func (m *Matrix)Add(that Matrix,  rowStart, colStart, rowCount int) {
	if(colStart + m.Cols() > that.Cols()) {
		panic("out of range");
	}
	for i := 0; i < m.Cols(); i++ {
		if(that.ValueCount(colStart + i) != m.ValueCount(i)){
			panic("incompatible relations");
		}
	}
	for j := 0; j < rowCount; j++ {
		rowSrc := that.Row(rowStart + j);
		rowDest := make([]float64,m.Cols());
		for i := 0;	i < m.Cols();	i++{
			rowDest[i] = rowSrc[colStart + i];
		}
		m.m_data = append(m.m_data,rowDest);
	}

}

// Resizes this matrix (and sets all attributes to be continuous)
func (m *Matrix)SetSize(rows,cols int) {
	for j := 0; j < rows; j++ {
		row := make([]float64,cols);
		m.m_data = append(m.m_data,row);
	}
	for i := 0; i < cols; i++ {
		m.m_attr_name = append(m.m_attr_name,"");
		m.m_str_to_enum = append(m.m_str_to_enum,treemap.NewWithStringComparator());
		m.m_enum_to_str = append(m.m_enum_to_str,treemap.NewWithIntComparator());
	}
}

// Loads from an ARFF file
func (m *Matrix)LoadArff(filename string) {
 	READDATA := false;
 	f, _:= os.Open(filename)
 	src := bufio.NewScanner(f);
	// Initialize the scanner.
	src.Split(bufio.ScanLines)
	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		success := src.Scan()
		if !success {
			if src.Err() == nil {
				break;
			} else {
				panic(src.Err());
			}
		} else {
			line := src.Text();
			tokens := strings.Fields(line);
			fmt.Println("Tokens ",tokens);
			if len(tokens) > 0 && tokens[0] != "" && tokens[0][0] != '%'{
				if !READDATA{
					firstToken := strings.ToUpper(tokens[0]);
					fmt.Println("First Token ", firstToken)
					if firstToken == "@RELATION"{
						datasetName := tokens[1];
						fmt.Println("Dataset Name: ",datasetName);
					}
					if firstToken == "@ATTRIBUTE"{
						ste := treemap.NewWithStringComparator();
						m.m_str_to_enum = append(m.m_str_to_enum,ste);
						ets := treemap.NewWithIntComparator();
						m.m_enum_to_str = append(m.m_enum_to_str,ets);
						var attributeName string;
						if strings.Index(line,"'") != -1{
							tokens := strings.Split(line,"'");
							attributeName = tokens[1];
						}else {
							attributeName = tokens[1];
						}
						m.m_attr_name = append(m.m_attr_name,attributeName);
						vals := 0;
						typ := strings.ToUpper(tokens[2]);
						if typ == "REAL" || typ == "CONTINUOUS" || typ == "INTEGER"{

						}else{
							values := strings.Split(line[strings.Index(line,"{")+1:strings.Index(line,"}")],",")
							fmt.Println("Values: ", values);
							for index := range values {
								ste.Put(values[index],vals);
								ets.Put(vals,values[index]);
								vals++;
							}
						}
					}
					if firstToken == "@DATA"{
						READDATA = true;
					}
				}else{
					newrow := make([]float64,m.Cols())
					curPos := 0;
					values := strings.Split(line,",")
					for index := range values {
						var doubleValue float64;
						vals := m.m_enum_to_str[curPos].Size();
						if values[index] == "?"{
							doubleValue = MISSING;
						}else if vals == 0{
							doubleValue , _= strconv.ParseFloat(values[index],64);
						}else{
							val1, _ := m.m_str_to_enum[curPos].Get(values[index])
							doubleValue = float64(val1.(int));
						}
						newrow[curPos] = doubleValue;
						curPos++;
					}
					m.m_data = append(m.m_data,newrow);
				}
			}
		}
	}
}


// Returns the number of rows in the matrix
func (m *Matrix)Rows()int{
	return len(m.m_data);
}
// Returns the number of columns (or attributes) in the matrix
func (m *Matrix)Cols()int{
	return len(m.m_attr_name);
}
// Returns the specified row
func (m *Matrix)Row(r int)[]float64{
	return m.m_data[r];
}
// Returns the element at the specified row and column
func (m *Matrix)Get(r,c int)float64{
	return m.m_data[r][c];
}
// Sets the value at the specified row and column
func (m *Matrix)Set(r,c int,v float64){
	m.m_data[r][c] = v;
}
// Returns the name of the specified attribute
func (m *Matrix)AttrName(c int)string{
	return m.m_attr_name[c];
}
// Set the name of the specified attribute
func (m *Matrix)SetAttrName(c int,s string){
	m.m_attr_name[c] = s;
}

// Returns the name of the specified value
func (m *Matrix)AttrValue(attr,val int)string{
	value, found := m.m_enum_to_str[attr].Get(val);
	if (found){
		return value.(string);
	} else{
		return "";
	}
}
// Returns the number of values associated with the specified attribute (or column)
// 0=continuous, 2=binary, 3=trinary, etc.
func (m *Matrix)ValueCount(c int)int{
	return m.m_enum_to_str[c].Size();
}

// Shuffles the row order
func (m *Matrix)Shuffle(rand Random){
	for n := m.Rows(); n > 0; n--{
		i := rand.NextInt(uint64(n));
		tmp := m.Row(n-1);
		m.m_data[n-1] = m.Row(int(i))
		m.m_data[i] = tmp;
	}
}

// Shuffles the row order with a buddy matrix
func (m *Matrix)ShuffleWithBuddy(rand Random, buddy Matrix ) {
	for n := m.Rows(); n > 0; n-- {
		i := rand.NextInt(uint64(n));
		tmp := m.Row(n - 1);
		m.m_data[n - 1] =  m.Row(int(i));
		m.m_data[i] = tmp
		tmp1 := buddy.Row(n - 1);
		buddy.m_data [n - 1] =  buddy.Row(int(i));
		buddy.m_data[i] = tmp1;
	}
}

// Returns the mean of the specified column
func (m *Matrix)ColumnMean(col int) float64 {
	var sum,count float64;
 	count = 0;
 	sum = 0;
	for i := 0; i < m.Rows(); i++ {
		v := m.Get(i, col);
		if(v != MISSING) {
			sum += v;
			count++;
		}
	}
	return sum / count;
}
// Returns the min value in the specified column
func (m *Matrix)ColumnMin(col int) float64 {
	min := MISSING;
	for i := 0; i < m.Rows(); i++ {
		v := m.Get(i, col);
		if(v != MISSING) {
			if (min == MISSING || v < min){
				min = v;
			}
		}
	}
	return min;
}
// Returns the max value in the specified column
func (m *Matrix)ColumnMax(col int) float64 {
	max := MISSING;
	for i := 0; i < m.Rows(); i++ {
		v := m.Get(i, col);
		if(v != MISSING) {
			if (max == MISSING || v > max){
				max = v;
			}
		}
	}
	return max;
}



// Returns the most common value in the specified column
func (m *Matrix)MostCommonValue(col int)float64 {
	tm := treemap.NewWith(utils.Float64Comparator);
	for i := 0; i < m.Rows(); i++ {
		v := m.Get(i, col);
		if (v != MISSING){
			count, _ := tm.Get(v);
			if (count == nil){
				tm.Put(v,1);
			}else{
				tm.Put(v,count.(int)+ 1);
			}
		}
	}
	maxCount := 0;
	val := MISSING;
	it := tm.Iterator();
	it.Begin()
	for {
		exists := it.Next()
		if exists {
			if it.Value().(int) > maxCount {
				maxCount = it.Value().(int);
				val = it.Key().(float64);
			}
		}
	}
	return val;
}



func (m *Matrix) Normalize() {
	for i := 0; i < m.Cols(); i++ {
		if (m.ValueCount(i) == 0) {
			min := m.ColumnMin(i);
			max := m.ColumnMax(i);
			for j := 0; j < m.Rows(); j++ {
				v := m.Get(j, i);
				if (v != MISSING) {
					m.Set(j, i, (v-min)/(max-min));
				}
			}
		}
	}
}

func (m *Matrix)Print() {
	fmt.Println("@RELATION Untitled");
	for i := 0; i < len(m.m_attr_name); i++ {
		fmt.Print("@ATTRIBUTE " + m.m_attr_name[i]);
		vals := m.ValueCount(i);
		if (vals == 0) {
			fmt.Println(" CONTINUOUS");
		} else {
			fmt.Print(" {");
			for j := 0; j < vals; j++ {
				if (j > 0) {
					fmt.Print(", ")
				}
				fmt.Print(m.m_enum_to_str[i].Get(j));
			}
			fmt.Println("}");
		}
	}
	fmt.Println("@DATA");
	for i := 0; i < m.Rows(); i++ {
		r := m.Row(i);
		for j := 0; j < len(r); j++ {
			if (j > 0) {
				fmt.Print(", ");
			}
			if (m.ValueCount(j) == 0) {
				fmt.Print(r[j]);
			} else {
				val, _ := m.m_enum_to_str[j].Get(int(r[j]))
				fmt.Print(val);
			}
		}
		fmt.Println("");

	}
}