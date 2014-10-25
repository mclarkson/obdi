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
mgrApp.directive("updsel", function () {
// ------------------------------------------------------------------------
// Refreshes the bootstrap-select add-on after angularjs has
// loaded dcs variable.

    return function (scope, element, attrs) {
        scope.$watch("dcs", function (value) {//I change here
            var val = value || null;            
            if (val)
                element.selectpicker('refresh');
                //$('.selectpicker').selectpicker('refresh');
        });
    };
});

// ------------------------------------------------------------------------
mgrApp.directive("updenvcapsel", ['$timeout', function ( $timeout ) {
// ------------------------------------------------------------------------
// Used in addenv-capstab.html (Capabilities tab)
// Refreshes the bootstrap-select add-on after angularjs has
// processed the template.

    return function (scope, element, attrs) {
      scope.$watch("newcap.selected", function (value) {
          $timeout( function() {
              element.selectpicker('val',element.val());
              element.selectpicker('refresh');
              //$('.selectpicker').selectpicker('refresh');
          });
      });
      scope.$watch("envcaps", function (value) {
          $timeout( function() {
              element.selectpicker('val',element.val());
              element.selectpicker('refresh');
              //$('.selectpicker').selectpicker('refresh');
          });
      });
    }
}]);

// ------------------------------------------------------------------------
mgrApp.controller("envCtrl", function ($log, $modal, $scope, $http, baseUrl) {
// ------------------------------------------------------------------------

  $scope.addenv = false;
  $scope.editenv = false;
  $scope.env = {};
  $scope.envs = {};
  $scope.envcapmaps = [];
  $scope.newcap = {};
  $scope.newcap.newcapmaps = [];
  $scope.newcap.selected = {};
  $scope.delcaps = {};
  $scope.delcaps.ids = [];
  $scope.dcs = {};
  $scope.managecaps = false;

  $scope.addenvtabs = [
    { title:'Environment Details', content:'frag/addenv-detailstab.html'},
    { title:'Capabilities', content:'frag/addenv-capstab.html'}
    ];

  $scope.editenvtabs = [
    { title:'Environment Details', content:'frag/editenv-detailstab.html'},
    { title:'Capabilities', content:'frag/editenv-capstab.html'}
    ];

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
      //$log.info('$scope.dc.Login: ' + $scope.dc.Login);

      // Set dc_obj so the select control works
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
      //$scope.FillEnvTable();

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
      if( found.length == 0 ) {
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

    var modalInstance = $modal.open({
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
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, SysName,
                                       DcSysName) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.SysName = SysName;
    $scope.DcSysName = DcSysName;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

  $scope.FillEnvTable();
  $scope.FillEnvCapsTable(); // For Capabilities dropdown

  // --------------------------------------------------------------------
  $scope.Edit_DeleteModal = function (id,Code) {
  // --------------------------------------------------------------------

    var modalInstance = $modal.open({
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

    var modalInstance = $modal.open({
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
  $scope.ModalDeleteInstanceCtrl = function ($scope, $modalInstance, Code) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.Code = Code;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

});

