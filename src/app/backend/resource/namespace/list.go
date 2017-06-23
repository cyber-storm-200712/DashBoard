// Copyright 2017 The Kubernetes Dashboard Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package namespace

import (
	"log"

	"github.com/kubernetes/dashboard/src/app/backend/api"
	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// NamespaceList contains a list of namespaces in the cluster.
type NamespaceList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of Namespaces.
	Namespaces []Namespace `json:"namespaces"`
}

// Namespace is a presentation layer view of Kubernetes namespaces. This means it is namespace plus
// additional augmented data we can get from other sources.
type Namespace struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Phase is the current lifecycle phase of the namespace.
	Phase v1.NamespacePhase `json:"phase"`
}

// GetNamespaceListFromChannels returns a list of all namespaces in the cluster.
func GetNamespaceListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (*NamespaceList,
	error) {
	log.Print("Getting namespace list")

	namespaces := <-channels.NamespaceList.List
	if err := <-channels.NamespaceList.Error; err != nil {
		return &NamespaceList{Namespaces: []Namespace{}}, err
	}

	return toNamespaceList(namespaces.Items, dsQuery), nil
}

// GetNamespaceList returns a list of all namespaces in the cluster.
func GetNamespaceList(client *client.Clientset, dsQuery *dataselect.DataSelectQuery) (*NamespaceList,
	error) {
	log.Print("Getting namespace list")

	namespaces, err := client.Namespaces().List(metaV1.ListOptions{
		LabelSelector: labels.Everything().String(),
		FieldSelector: fields.Everything().String(),
	})

	if err != nil {
		return nil, err
	}

	return toNamespaceList(namespaces.Items, dsQuery), nil
}

func toNamespaceList(namespaces []v1.Namespace, dsQuery *dataselect.DataSelectQuery) *NamespaceList {
	namespaceList := &NamespaceList{
		Namespaces: make([]Namespace, 0),
		ListMeta:   api.ListMeta{TotalItems: len(namespaces)},
	}

	namespaceCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(namespaces), dsQuery)
	namespaces = fromCells(namespaceCells)
	namespaceList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, namespace := range namespaces {
		namespaceList.Namespaces = append(namespaceList.Namespaces, toNamespace(namespace))
	}

	return namespaceList
}

func toNamespace(namespace v1.Namespace) Namespace {
	return Namespace{
		ObjectMeta: api.NewObjectMeta(namespace.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindNamespace),
		Phase:      namespace.Status.Phase,
	}
}
