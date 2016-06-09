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

import {StateParams} from 'common/resource/resourcedetail';
import {stateName} from 'replicationcontrollerdetail/replicationcontrollerdetail_state';

/**
 * Controller for the replication controller card menu
 *
 * @final
 */
export default class ReplicationControllerCardMenuController {
  /**
   * @param {!ui.router.$state} $state
   * @param
   * {!./../replicationcontrollerdetail/replicationcontroller_service.ReplicationControllerService}
   * kdReplicationControllerService
   * @ngInject
   */
  constructor($state, kdReplicationControllerService) {
    /**
     * Initialized from the scope.
     * @export {!backendApi.ReplicationController}
     */
    this.replicationController;

    /** @private {!ui.router.$state} */
    this.state_ = $state;

    /** @private
     * {!./../replicationcontrollerdetail/replicationcontroller_service.ReplicationControllerService}
     */
    this.kdReplicationControllerService_ = kdReplicationControllerService;

    /** @export */
    this.i18n = i18n;
  }

  /**
   * @param {!function(!MouseEvent)} $mdOpenMenu
   * @param {!MouseEvent} $event
   * @export
   */
  openMenu($mdOpenMenu, $event) { $mdOpenMenu($event); }

  /**
   * @export
   */
  viewDetails() {
    this.state_.go(
        stateName, new StateParams(
                       this.replicationController.objectMeta.namespace,
                       this.replicationController.objectMeta.name));
  }

  /**
   * @export
   */
  showUpdateReplicasDialog() {
    this.kdReplicationControllerService_.showUpdateReplicasDialog(
        this.replicationController.objectMeta.namespace, this.replicationController.objectMeta.name,
        this.replicationController.pods.current, this.replicationController.pods.desired);
  }
}

/**
 * @return {!angular.Component}
 */
export const replicationControllerCardMenuComponent = {
  bindings: {
    'replicationController': '=',
  },
  controller: ReplicationControllerCardMenuController,
  templateUrl: 'replicationcontrollerlist/replicationcontrollercardmenu.html',
};

const i18n = {
  /** @export {string} @desc Action 'View details' on the drop down menu for a
      single replication controller (replication controller list page).*/
  MSG_RC_LIST_VIEW_DETAILS_ACTION: goog.getMsg('View details'),
  /** @export {string} @desc Action 'Edit Pod Count' on the drop down menu for a
      single replication controller (replication controller list page).*/
  MSG_RC_LIST_EDIT_POD_COUNT_ACTION: goog.getMsg('Scale'),
  /** @export {string} @desc Label 'Replication Controller' which will appear in the replication Controller
      delete dialogm opened from a replication controller card on the list page. */
  MSG_RC_LIST_REPLICATION_CONTROLLER_LABEL: goog.getMsg('Replication Controller'),
};
