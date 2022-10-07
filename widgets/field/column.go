package field

import (
	"fmt"
	"strings"

	"github.com/yaoapp/yao/widgets/component"
	"github.com/yaoapp/yao/widgets/expression"
)

// Replace replace with data
func (column ColumnDSL) Replace(data map[string]interface{}) (*ColumnDSL, error) {
	new := column
	err := expression.Replace(&new.Key, data)
	if err != nil {
		return nil, err
	}

	err = expression.Replace(&new.Bind, data)
	if err != nil {
		return nil, err
	}

	if new.Edit != nil {
		err = expression.Replace(&new.Edit.Props, data)
		if err != nil {
			return nil, err
		}
	}

	if new.View != nil {
		err = expression.Replace(&new.View.Props, data)
		if err != nil {
			return nil, err
		}
	}

	return &new, nil
}

// Trans trans
func (column *ColumnDSL) Trans(widgetName string, inst string, trans func(widget string, inst string, value *string) bool) bool {
	res := false
	if column.Edit != nil {
		if column.Edit.Trans(widgetName, inst, trans) {
			res = true
		}
	}

	if column.View != nil {
		if column.View.Trans(widgetName, inst, trans) {
			res = true
		}
	}

	return res
}

// Trans column trans
func (columns Columns) Trans(widgetName string, inst string, trans func(widget string, inst string, value *string) bool) bool {
	res := false

	for key, column := range columns {
		if trans(widgetName, inst, &key) {
			res = true
		}
		newPtr := &column
		if newPtr.Trans(widgetName, inst, trans) {
			res = true
		}
		columns[key] = *newPtr
	}

	return res
}

// Map cast to map[string]inteface{}
func (column ColumnDSL) Map() map[string]interface{} {
	res := map[string]interface{}{
		"key":  column.Key,
		"link": column.Link,
		"bind": column.Bind,
	}

	if column.View != nil {
		res["view"] = column.View.Map()
	}

	if column.Edit != nil {
		res["edit"] = column.Edit.Map()
	}
	return res
}

// CPropsMerge merge the Columns cloud props
func (columns Columns) CPropsMerge(cloudProps map[string]component.CloudPropsDSL, getXpath func(name string, kind string, column ColumnDSL) (xpath string)) error {

	for name, column := range columns {

		if column.Edit != nil && column.Edit.Props != nil {
			xpath := getXpath(name, "edit", column)
			cProps, err := column.Edit.Props.CloudProps(xpath)
			if err != nil {
				return err
			}
			mergeCProps(cloudProps, cProps)
		}

		if column.View != nil && column.View.Props != nil {
			xpath := getXpath(name, "view", column)
			cProps, err := column.View.Props.CloudProps(xpath)
			if err != nil {
				return err
			}
			mergeCProps(cloudProps, cProps)
		}
	}

	return nil
}

// ComputeFieldsMerge merge the compute fields
func (columns Columns) ComputeFieldsMerge(computeInFields map[string]string, computeOutFields map[string]string) {
	for name, column := range columns {

		// Compute In
		if column.In != "" {
			if !strings.Contains(column.In, ".") {
				column.In = fmt.Sprintf("yao.component.%s", column.In)
			}
			computeInFields[column.Bind] = column.In
			computeInFields[name] = column.In
		}

		// Compute Out
		if column.Out != "" {
			if !strings.Contains(column.Out, ".") {
				column.In = fmt.Sprintf("yao.component.%s", column.Out)
			}
			computeOutFields[column.Bind] = column.Out
			computeOutFields[name] = column.Out
		}
	}
}

func mergeCProps(cloudProps map[string]component.CloudPropsDSL, cProps map[string]component.CloudPropsDSL) {
	for k, v := range cProps {
		cloudProps[k] = v
	}
}
