// Copyright 2016 Mark Clarkson
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

// ------------------------------------------------------------------------
mgrApp.controller("capsCtrl", function ($log, $uibModal, $scope, $http, baseUrl) {
// ------------------------------------------------------------------------

  $scope.adddc = false;
  $scope.editdc = false;
  $scope.dccaps = {};
  $scope.newcap = {};

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.mainmessage = "";
    $scope.mainokmessage = "";
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    clearMessages();

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccaps/" + id
    }).success( function(data, status, headers, config) {
      //$scope.mainokmessage = "The dc was deleted."
      $scope.FillCapsTable();
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server error: " + data['Error'];
        $scope.error = true;
      } else if (status==0) {
        $scope.errtext = "Could not connect to server.";
        $scope.error = true;
      } else {
        $scope.errtext = "Unknown error.";
        $scope.error = true;
      }
    });
  }

  // ----------------------------------------------------------------------
  $scope.AddCapability = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new user

    clearMessages();

    $http({
      method: 'POST',
      data: $scope.newcap,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccaps"
    }).success( function(data, status, headers, config) {
      //$scope.okmessage = "The dc was added."
      $scope.FillCapsTable();
      $scope.newcap = {};
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server error: " + data['Error'];
        $scope.error = true;
      } else if (status==0) {
        $scope.errtext = "Could not connect to server.";
        $scope.error = true;
      } else {
        $scope.errtext = "Unknown error.";
        $scope.error = true;
      }
    });
  }

  // ----------------------------------------------------------------------
  $scope.FillCapsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccaps"
    }).success( function(data, status, headers, config) {
      $scope.dccaps = data;
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==0) {
        // This is a guess really
        $scope.login.errtext = "Could not connect to server.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else {
        $scope.login.errtext = "Logged out due to an unknown error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      }
    });
  };

  $scope.FillCapsTable();

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.dialog = function (id,Code) {
  // --------------------------------------------------------------------

    $scope.Code = Code;
    $scope.id = id;

    var modalInstance = $uibModal.open({
      templateUrl: 'DeleteDccap.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        Code: function () {
          return $scope.Code;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.Code + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // ----------------------------------------------------------------------
  $scope.UpdateCap = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new user

    clearMessages();

    $http({
      method: 'PUT',
      data: $scope.dccap,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccaps/" + $scope.dccap.Id
    }).success( function(data, status, headers, config) {
      //$scope.okmessage = "The dc was added."
      $scope.FillCapsTable();
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server error: " + data['Error'];
        $scope.error = true;
      } else if (status==0) {
        $scope.errtext = "Could not connect to server.";
        $scope.error = true;
      } else {
        $scope.errtext = "Unknown error.";
        $scope.error = true;
      }
    });
  }

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $uibModalInstance, Code) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.Code = Code;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.editdialog = function (dccap) {
  // --------------------------------------------------------------------

    $scope.dccap = dccap;
    //$scope.id = id;

    var modalInstance = $uibModal.open({
      templateUrl: 'EditDccap.html',
      controller: $scope.ModalInstanceEditCtrl,
      size: 'md',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        dccap: function () {
          return $scope.dccap;
        }
      }
    });

    modalInstance.result.then(function () {
      $log.info('Will update: ' + $scope.dccap.Code );
      $scope.UpdateCap($scope.dccap);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceEditCtrl = function ($scope, $uibModalInstance, dccap) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.dccap = dccap;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

});

