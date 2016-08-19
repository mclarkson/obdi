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
mgrApp.controller("envCtrl", function ($log, $uibModal, $scope, $http,
baseUrl, $timeout) {
// ------------------------------------------------------------------------

  $scope.addenv = false;
  $scope.editenv = false;
  $scope.env = {};
  $scope.env.SysName = "";
  $scope.env.WorkerUrl = "";
  $scope.envs = {};
  $scope.envcapmaps = [];
  $scope.newcap = {};
  $scope.newcap.newcapmaps = [];
  $scope.newcap.selected = {};
  $scope.delcaps = {};
  $scope.delcaps.ids = [];
  $scope.dcs = {};
  $scope.managecaps = false;
  $scope.worker = {};
  $scope.jsonobject = {};

  $scope.addenvtabs = [
    { title:'Environment Details', content:'frag/addenv-detailstab.html'},
    { title:'Capabilities', content:'frag/addenv-capstab.html'}
    ];

  $scope.editenvtabs = [
    { title:'Environment Details', content:'frag/editenv-detailstab.html'},
    { title:'Capabilities', content:'frag/editenv-capstab.html'}
    ];

  // ----------------------------------------------------------------------
  $scope.$watch('env.SysName', function() {
  // ----------------------------------------------------------------------
    $scope.env.SysName = $scope.env.SysName.replace(/\s+/g,'');
  });

  // ----------------------------------------------------------------------
  $scope.$watch('env.WorkerUrl', function() {
  // ----------------------------------------------------------------------
    $scope.env.WorkerUrl = $scope.env.WorkerUrl.replace(/\s+/g,'');
  });

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.mainmessage = "";
    $scope.mainokmessage = "";
  }

  // ----------------------------------------------------------------------
  $scope.AddEnv = function(tf) {
  // ----------------------------------------------------------------------
    $scope.addenv = tf;
    $scope.env = {};
    $scope.env.SysName = "";
    $scope.env.WorkerUrl = "";
    $scope.newcap = {};
    $scope.newcap.newcapmaps = [];
    $scope.newcap.selected = {};
    $scope.FillDcTable();
    $scope.FillEnvCapsTable();
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.EditEnv = function(id) {
  // ----------------------------------------------------------------------
    $scope.editenv = id;

    if (id) {
      // Find the item id in the environments array
      $scope.env = $.grep($scope.envs, function(e){ return e.Id == id; })[0];

      // Set dc_obj since select is bound to it.
      // dc_obj is a copy so must use 'track by' in the view.
      $scope.env.dc_obj = $.grep($scope.dcs,function(e)
        {return e.Id == $scope.env.DcId; })[0];

      $scope.FillEnvCapsTable(); // For Capabilities
      $scope.FillEnvCapsMapsTable( id ); // For Capability Maps
      $scope.newcap.newcapmaps = [];
      $scope.newcap.selected = {};
    }

    $scope.FillEnvTable();

    clearMessages();

  }

  // ----------------------------------------------------------------------
  $scope.ManageCaps = function(tf) {
  // ----------------------------------------------------------------------
    $scope.managecaps = tf;
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs/" + id
    }).success( function(data, status, headers, config) {
      $scope.mainokmessage = "The environment was deleted."
      $scope.FillEnvTable();
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
  function AddCapMap_noasync( json_obj, type, id ) {
  // ----------------------------------------------------------------------

    json = JSON.stringify( json_obj )
    jQuery.ajax({
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envcapmaps" + id,
      data: json,
      type: type,
      processData: false,
      async:   false,
      error: function(jqXHR, textStatus, errorThrown) {
        var status = jqXHR.status;
        if (status>=500) {
          $scope.errtext = "Internal server error.";
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
      }
    });
  }

  // ----------------------------------------------------------------------
  $scope.EditApply = function() {
  // ----------------------------------------------------------------------
  // Apply button for editing a new user

    clearMessages();

    $scope.env.DcId = $scope.env.dc_obj.Id;

    $http({
      method: 'PUT',
      data: $scope.env,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs/" + $scope.env.Id
    }).success( function(data, status, headers, config) {

      // Cycle through each entry in $scope.newcap.newcapmap array
      // and REST ADD to /envcapmaps

      for( var i = 0; i < $scope.newcap.newcapmaps.length; i++ ) {
        json = {
                "EnvId":$scope.editenv,
                "EnvCapId":$scope.newcap.newcapmaps[i].Id
        };
        AddCapMap_noasync( json, 'POST', '' );
      }

      // Cycle through each entry in $scope.delcap.ids array
      // and REST DELETE to /envcapmaps

      for( var i = 0; i < $scope.delcaps.ids.length; i++ ) {
        json= {};
        AddCapMap_noasync( json, 'DELETE', '/'+$scope.delcaps.ids[i] );
      }

      $scope.delcaps = {};
      $scope.delcaps.ids = [];

      $scope.okmessage = "Changes were applied."
      $scope.FillEnvTable();
      $scope.FillEnvCapsTable(); // For Capabilities
      $scope.FillEnvCapsMapsTable($scope.env.Id);
      $scope.newcap.newcapmaps = [];
      $scope.newcap.selected = {};

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
  $scope.Apply = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new environment

    clearMessages();

    $scope.env.DcId = $scope.env.dc_obj.Id;
    //$log.info('$scope.dc_obj.Id: ' + $scope.dc_obj.Id);
    //alert($scope.env.dc_obj);
    //$log.info('$scope.env.DcId: ' + $scope.env.DcId);

    $http({
      method: 'POST',
      data: $scope.env,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs"
    }).success( function(data, status, headers, config) {

      id = data.Id;

      // Cycle through each entry in $scope.newcap.newcapmap array
      // and REST add to /envcapmaps

      for( var i = 0; i < $scope.newcap.newcapmaps.length; i++ ) {
        json = {
                "EnvId":id,
                "EnvCapId":$scope.newcap.newcapmaps[i].Id
        };
        AddCapMap_noasync( json, 'POST', '' );
      }

      $scope.okmessage = "The environment was added."
      $scope.FillEnvTable();

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
  $scope.FillEnvCapsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envcaps"
    }).success( function(data, status, headers, config) {
      $scope.envcaps = data;
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

  // ----------------------------------------------------------------------
  $scope.FillEnvCapsMapsTable = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envcapmaps?env_id=" + id
    }).success( function(data, status, headers, config) {
      $scope.envcapmaps = data;
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

  // ----------------------------------------------------------------------
  $scope.FillEnvTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs"
    }).success( function(data, status, headers, config) {
      $scope.envs = data;
      $scope.FillDcTable();
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

  // ----------------------------------------------------------------------
  $scope.FillDcTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dcs"
    }).success( function(data, status, headers, config) {
      $scope.dcs = data;
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

  // ----------------------------------------------------------------------
  $scope.EditAddCapMap = function() {
  // ----------------------------------------------------------------------
    $scope.delcaps = {};
    $scope.delcaps.ids = [];

    if( typeof $scope.newcap.selected !== 'undefined' &&
        $scope.newcap.selected ) {

      // Search for the item being added

      var found = $.grep( $scope.envcapmaps,
        function(e){ return e.EnvCapId == $scope.newcap.selected.Id; });

      // Don't add it if it's there
      if( found.length == 0 && $scope.newcap.selected.Code ) {
        $scope.newcap.newcapmaps.push($scope.newcap.selected);
        table_item = {
                      "EnvCapCode":$scope.newcap.selected.Code,
                      "EnvCapDesc":$scope.newcap.selected.Desc,
                      "EnvCapId":$scope.newcap.selected.Id,
                      "EnvId":$scope.editenv
        };

        $scope.envcapmaps.push( table_item );
        clearMessages();
      }

      $scope.newcap.selected = {};
    }
  }

  // ----------------------------------------------------------------------
  $scope.AddCapMap = function() {
  // ----------------------------------------------------------------------
    if( typeof $scope.newcap.selected !== 'undefined' &&
        $scope.newcap.selected ) {

      // Search for the item being added
      var found = $.grep( $scope.newcap.newcapmaps,
        function(e){ return e.Id == $scope.newcap.selected.Id; });

      // Don't add it if was added already. (and don't report anything!)
      if( found.length == 0 ) {
        $scope.newcap.newcapmaps.push($scope.newcap.selected);
        clearMessages();
      }

      $scope.newcap.selected = {};
    }
  }

  // ----------------------------------------------------------------------
  $scope.Add_DeleteCapMapFromTable = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    // Delete from the display table
    $scope.newcap.newcapmaps = $.grep( $scope.newcap.newcapmaps,
      function(e){ return e.Id != id; });
  }

  // ----------------------------------------------------------------------
  $scope.Edit_DeleteCapMapFromTable = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    // Delete from the display table
    $scope.envcapmaps = $.grep( $scope.envcapmaps,
      function(e){ return e.Id != id; });

    // Delete from the newcap table in case it was just added
    $scope.newcap.newcapmaps = $.grep( $scope.newcap.newcapmaps,
      function(e){ return e.EnvCapId != id; });

    // Add to the delete-later array
    if( typeof id != 'undefined' && id ) {
      $scope.delcaps.ids.push( id );
    }
  }

  // ----------------------------------------------------------------------
  $scope.AddJsonObjectEntry = function( newJsonObject ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'POST',
      data: newJsonObject,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/jsonobjects"
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "The json object details were updated."
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // ----------------------------------------------------------------------
  $scope.UpdateJsonObjectEntry = function( newJsonObject, id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'PUT',
      data: newJsonObject,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/jsonobjects/" + id
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "The json object details were updated."
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // ----------------------------------------------------------------------
  $scope.AddWorkerEntry = function( newWorker ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'POST',
      data: newWorker,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/workers"
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "The worker details were updated."
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // ----------------------------------------------------------------------
  $scope.UpdateWorkerEntry = function( newWorker, id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'PUT',
      data: newWorker,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/workers/" + id
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "The worker details were updated."
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // ----------------------------------------------------------------------
  $scope.FillEnvCapsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envcaps"
    }).success( function(data, status, headers, config) {
      $scope.envcaps = data;
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

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.dialog = function (id,SysName,DcSysName) {
  // --------------------------------------------------------------------

    $scope.SysName = SysName;
    $scope.DcSysName = DcSysName;
    $scope.id = id;

    var modalInstance = $uibModal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // these variables are passed to the ModalInstanceCtrl
        SysName: function () {
          return $scope.SysName;
        },
        DcSysName: function () {
          return $scope.DcSysName;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.SysName + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $uibModalInstance, SysName,
                                       DcSysName) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.SysName = SysName;
    $scope.DcSysName = DcSysName;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.Edit_DeleteModal = function (id,Code) {
  // --------------------------------------------------------------------

    var modalInstance = $uibModal.open({
      templateUrl: 'DeleteEnvcapmap.html',
      controller: $scope.ModalDeleteInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        Code: function () {
          return Code;
        }
      }
    });

    modalInstance.result.then(function () {
      $log.info('Will delete: ' + Code + '(' + id + ')' );
      $scope.Edit_DeleteCapMapFromTable(id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.Add_DeleteModal = function (id,Code) {
  // --------------------------------------------------------------------

    var modalInstance = $uibModal.open({
      templateUrl: 'DeleteEnvcapmap.html',
      controller: $scope.ModalDeleteInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        Code: function () {
          return Code;
        }
      }
    });

    modalInstance.result.then(function () {
      $log.info('Will delete: ' + Code + '(' + id + ')' );
      $scope.Add_DeleteCapMapFromTable(id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalDeleteInstanceCtrl = function ($scope, $uibModalInstance, Code) {
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

  // ====================================================================

  // --------------------------------------------------------------------
  $scope.Edit_JsonObjectModal = function (EnvId, EnvCapId, EnvCapCode) {
  // --------------------------------------------------------------------

    $scope.Json = "";
    $scope.id = EnvCapId;

    var recordexists=false;

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/jsonobjects?env_id=" + EnvId + "&env_cap_id=" + EnvCapId
    }).success( function(data, status, headers, config) {
      $scope.jsonobjects = data;
      if( $scope.jsonobjects[0] ) recordexists=true;
      var modalInstance = $uibModal.open({
        templateUrl: 'EditJsonObjectDef.html',
        controller: $scope.Edit_JsonObjectModalCtrl,
        size: 'md',
        resolve: {
          // these variables are passed to the ModalInstanceCtrl
          EnvCapCode: function () {
            return EnvCapCode;
          },
          Json: function () {
            if( recordexists ) {
                return $scope.jsonobjects[0].Json;
            } else {
                return "";
            }
          },
        }
      });

      modalInstance.result.then(function (result) {

        var newJsonObject = {};
        newJsonObject.EnvId = EnvId;
        newJsonObject.EnvCapId = EnvCapId;
        newJsonObject.Json = result.Json;

        if( recordexists ) {
            return $scope.UpdateJsonObjectEntry(newJsonObject,$scope.jsonobjects[0].Id);
        } else {
            return $scope.AddJsonObjectEntry(newJsonObject);
        }
      }, function () {
        $log.info('Modal dismissed at: ' + new Date());
      });
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

  // --------------------------------------------------------------------
  $scope.Edit_JsonObjectModalCtrl = function ($scope, $uibModalInstance,
                                EnvCapCode, Json) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.EnvCapCode = EnvCapCode;
    $scope.Json = Json;

    $scope.ok = function () {
      result = {};
      result.Json = $scope.Json;
      $uibModalInstance.close(result);
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.Edit_WorkerModal = function (EnvId, EnvCapId, EnvCapCode) {
  // --------------------------------------------------------------------

    $scope.WorkerUrl = "";
    $scope.WorkerKey = "";
    $scope.id = EnvCapId;

    var recordexists=false;

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/workers?env_id=" + EnvId + "&env_cap_id=" + EnvCapId
    }).success( function(data, status, headers, config) {
      $scope.worker = data;
      if( $scope.worker[0] ) recordexists=true;
      var modalInstance = $uibModal.open({
        templateUrl: 'EditWorkerDef.html',
        controller: $scope.Edit_WorkerModalCtrl,
        size: 'md',
        resolve: {
          // these variables are passed to the ModalInstanceCtrl
          EnvCapCode: function () {
            return EnvCapCode;
          },
          WorkerUrl: function () {
            if( recordexists ) {
                return $scope.worker[0].WorkerUrl;
            } else {
                return "";
            }
          },
          WorkerKey: function () {
            if( recordexists ) {
                return $scope.worker[0].WorkerKey;
            } else {
                return "";
            }
          },
        }
      });

      modalInstance.result.then(function (result) {

        var newWorker = {};
        newWorker.EnvId = EnvId;
        newWorker.EnvCapId = EnvCapId;
        newWorker.WorkerUrl = result.WorkerUrl;
        newWorker.WorkerKey = result.WorkerKey;

        if( recordexists ) {
            return $scope.UpdateWorkerEntry(newWorker,$scope.worker[0].Id);
        } else {
            return $scope.AddWorkerEntry(newWorker);
        }
      }, function () {
        $log.info('Modal dismissed at: ' + new Date());
      });
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

  // --------------------------------------------------------------------
  $scope.Edit_WorkerModalCtrl = function ($scope, $uibModalInstance,
                                EnvCapCode, WorkerUrl, WorkerKey) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.EnvCapCode = EnvCapCode;
    $scope.WorkerUrl = WorkerUrl;
    $scope.WorkerKey = WorkerKey;

    $scope.ok = function () {
      result = {};
      result.WorkerUrl = $scope.WorkerUrl;
      result.WorkerKey = $scope.WorkerKey;
      $uibModalInstance.close(result);
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

  $scope.FillEnvTable();
  $scope.FillEnvCapsTable(); // For Capabilities dropdown

});

