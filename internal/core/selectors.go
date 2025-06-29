/*
 * Copyright 2025 Michele Zanotti <m.zanotti019@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func NewDuplicableListOptions() metav1.ListOptions {
	selector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      LABEL_DUPLICATED,
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		},
	}
	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&selector),
	}
}

func NewDuplicatedListOptions() metav1.ListOptions {
	selector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      LABEL_DUPLICATED,
				Operator: metav1.LabelSelectorOpExists,
			},
		},
	}
	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&selector),
	}
}
