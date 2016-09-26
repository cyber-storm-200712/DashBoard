// Copyright 2015 Google Inc. All Rights Reserved.
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

package replicationcontroller

import (
	client "k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/service"
)

// GetReplicationControllerServices returns list of services that are related to replication
// controller targeted by given name.
func GetReplicationControllerServices(client client.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, rcName string) (*service.ServiceList, error) {

	replicationController, err := client.ReplicationControllers(namespace).Get(rcName)
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		ServiceList: common.GetServiceListChannel(client, common.NewSameNamespaceQuery(namespace),
			1),
	}

	services := <-channels.ServiceList.List
	if err := <-channels.ServiceList.Error; err != nil {
		return nil, err
	}

	matchingServices := common.FilterNamespacedServicesBySelector(services.Items, namespace,
		replicationController.Spec.Selector)
	return service.CreateServiceList(matchingServices, dsQuery), nil
}
