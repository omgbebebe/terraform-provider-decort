/*
Copyright (c) 2019-2022 Digital Energy Cloud Solutions LLC. All Rights Reserved.
Authors:
Petr Krutov, <petr.krutov@digitalenergy.online>
Stanislav Solovev, <spsolovev@digitalenergy.online>
Kasim Baybikov, <kmbaybikov@basistech.ru>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flattens

import "strconv"

func FlattenMeta(m []interface{}) []string {
	output := []string{}
	for _, item := range m {
		switch d := item.(type) {
		case string:
			output = append(output, d)
		case int:
			output = append(output, strconv.Itoa(d))
		case int64:
			output = append(output, strconv.FormatInt(d, 10))
		case float64:
			output = append(output, strconv.FormatInt(int64(d), 10))
		default:
			output = append(output, "")
		}
	}
	return output
}
