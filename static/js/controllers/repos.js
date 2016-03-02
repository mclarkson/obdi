// Obdi - a REST interface and GUI for deploying software
// Copyright (C) 2014  Mark Clarkson
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ------------------------------------------------------------------------
mgrApp.controller("reposCtrl", function ($log, $uibModal, $scope, $http,
      baseUrl) {
// ------------------------------------------------------------------------

  $scope.adddc = false;
  $scope.editdc = false;
  $scope.repos = {};
  $scope.repo = {};

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
           + "/repos/" + id
    }).success( function(data, status, headers, config) {
      //$scope.mainokmessage = "The dc was deleted."
      $scope.FillReposTable();
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
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
  $scope.AddRepo = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new user

    clearMessages();

    $http({
      method: 'POST',
      data: $scope.repo,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/repos"
    }).success( function(data, status, headers, config) {
      $scope.FillReposTable();
      $scope.repo = {};
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
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
  $scope.FillReposTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/repos"
    }).success( function(data, status, headers, config) {
      $scope.repos = data;
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

  $scope.FillReposTable();

  // Modal dialog

  // DELETE

  // --------------------------------------------------------------------
  $scope.dialog = function (id,Url) {
  // --------------------------------------------------------------------

    $scope.Url = Url;
    $scope.id = id;

    var modalInstance = $uibModal.open({
      templateUrl: 'DeleteRepo.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        Url: function () {
          return $scope.Url;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.Url + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // EDIT

  // ----------------------------------------------------------------------
  $scope.UpdateRepo = function() {
  // ----------------------------------------------------------------------
  // UNUSED

    clearMessages();

    $http({
      method: 'PUT',
      data: $scope.reponew,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/repos/" + $scope.reponew.Id
    }).success( function(data, status, headers, config) {
      //$scope.okmessage = "The dc was added."
      $scope.FillReposTable();
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server error.";
        $scope.error = true;
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
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
  $scope.ModalInstanceCtrl = function ($scope, $uibModalInstance, Url) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.Url = Url;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.editdialog = function (repo) {
  // --------------------------------------------------------------------
  // UNUSED

    $scope.reponew = {};
    $scope.reponew.Id = repo.Id;
    $scope.reponew.Url = repo.Url;

    var modalInstance = $uibModal.open({
      templateUrl: 'EditRepo.html',
      controller: $scope.ModalInstanceEditCtrl,
      size: 'md',
      resolve: {
        reponew: function () {
          return $scope.reponew;
        }
      }
    });

    modalInstance.result.then(function () {
      $log.info('Will update: ' + $scope.reponew.Url );
      $scope.UpdateRepo($scope.reponew);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceEditCtrl = function ($scope, $uibModalInstance, reponew) {
  // --------------------------------------------------------------------
  // UNUSED

    // So the template can access 'loginname' in this new scope
    $scope.reponew = reponew;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

});

