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

// // ------------------------------------------------------------------------
// mgrApp.directive("uppluginapsel", ['$timeout', function ( $timeout ) {
// // ------------------------------------------------------------------------
// // Refreshes the bootstrap-select add-on after angularjs has
// // processed the template.
// 
//     return function (scope, element, attrs) {
//       scope.$watch("newcap.selected", function (value) {
//           $timeout( function() {
//               element.selectpicker('val',element.val());
//               element.selectpicker('refresh');
//               //$('.selectpicker').selectpicker('refresh');
//           });
//       });
//       scope.$watch("plugincaps", function (value) {
//           $timeout( function() {
//               element.selectpicker('val',element.val());
//               element.selectpicker('refresh');
//               //$('.selectpicker').selectpicker('refresh');
//           });
//       });
//     }
// }]);

// ------------------------------------------------------------------------
mgrApp.controller("pluginCtrl", function ($log, $modal, $scope, $http,
            baseUrl) {
// ------------------------------------------------------------------------

  $scope.addplugins = false;
  $scope.editplugins = false;
  $scope.plugincapmaps = [];
  $scope.plugincaps = [];
  $scope.newcap = {};
  $scope.newcap.newcapmaps = [];
  $scope.newcap.selected = {};
  $scope.plugins = {};
  $scope.pluginsAvail = [];
  $scope.pluginsAvailFinished = false;
  $scope.delcaps = {};
  $scope.delcaps.ids = [];
  $scope.managerepos = false;

  //$scope.addplugintabs = [];
    //{ title:'Plugin Details', content:'frag/addplugin-detailstab.html'},
    //{ title:'Files', content:'frag/addplugin-files.html'}
    //];

  $scope.editplugintabs = [
    { title:'Plugin Details', content:'frag/editplugin-detailstab.html'},
    { title:'Files', content:'frag/editplugin-files.html'}
    ];

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.mainmessage = "";
    $scope.mainokmessage = "";
    $scope.login.error = false;
    $scope.error = false;
  }

  // ----------------------------------------------------------------------
  $scope.AddPlugin = function(tf) {
  // ----------------------------------------------------------------------
    $scope.addplugins = tf;
    $scope.pluginsAvailMeta = [];
    if(tf) {
      $scope.pluginsAvail = [];
      $scope.pluginsAvailFinished = false;
      $scope.FillPluginsAvailTable(-1);
    } else {
      $scope.FillPluginTable();
    }
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.RestAddPlugin = function(name) {
  // ----------------------------------------------------------------------

    clearMessages();

    var index = -1;
    for( var i=0; i<$scope.pluginsAvail.length; ++i ) {
      if( $scope.pluginsAvail[i].Name == name ) {
        index = i;
        break;
      }
    }
    $scope.pluginsAvailMeta[index] = {};
    $scope.pluginsAvailMeta[index].Completed = false;
    $scope.pluginsAvailMeta[index].Disabled = true;
    var data = {};
    data.Name = $scope.pluginsAvail[index].Name;

    $http({
      method: 'POST',
      data: data,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/repoplugins"
    }).success( function(data, status, headers, config) {
      //$scope.FillPluginsAvailTable(index);
      var index = -1;
      // Delete from the meta table
      for( var i=0; i<$scope.pluginsAvail.length; ++i ) {
        if( $scope.pluginsAvail[i].Name == name ) {
          index = i;
          break;
        }
      }
      $scope.pluginsAvailMeta.splice(index,1);
      $scope.pluginsAvail.splice(index,1);
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
  }

  // ----------------------------------------------------------------------
  $scope.EditPlugin = function(id) {
  // ----------------------------------------------------------------------
    $scope.editplugins = id;
    $scope.delcaps = {};
    $scope.delcaps.ids = [];

    if (id) {
      $scope.plugin = $.grep($scope.plugins, function(e){ return e.Id == id; })[0];
      $scope.newcap.newcapmaps = [];
      $scope.newcap.selected = {};
    }

    $scope.FillPluginTable();
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.ManageRepos = function(tf) {
  // ----------------------------------------------------------------------
    $scope.managerepos = tf;
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/plugins/" + id
    }).success( function(data, status, headers, config) {
      $scope.mainokmessage = "The plugin was deleted."
      $scope.FillPluginTable();
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
    $scope.plugincapmaps = $.grep( $scope.plugincapmaps,
      function(e){ return e.Id != id; });

    // Delete from the newcap table in case it was just added
    $scope.newcap.newcapmaps = $.grep( $scope.newcap.newcapmaps,
      function(e){ return e.PluginCapId != id; });

    // Add to the delete-later array
    if( typeof id != 'undefined' && id ) {
      $scope.delcaps.ids.push( id );
    }
  }

  // ----------------------------------------------------------------------
  $scope.EditApply = function() {
  // ----------------------------------------------------------------------
  // Apply button for editing a new user

    clearMessages();

    $http({
      method: 'PUT',
      data: $scope.plugin,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/plugins/" + $scope.plugin.Id
    }).success( function(data, status, headers, config) {

      $scope.okmessage = "Changes were applied."
      $scope.FillPluginTable();

    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server said.";
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
           + "/plugincapmaps" + id,
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
  $scope.Apply = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new user

    $scope.message = "";

    $http({
      method: 'POST',
      data: $scope.plugin,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/plugins"
    }).success( function(data, status, headers, config) {

      id = data.Id;

      // Cycle through each entry in $scope.newcap.newcapmap array
      // and REST add to /envcapmaps

      for( var i = 0; i < $scope.newcap.newcapmaps.length; i++ ) {
        json = {
                "PluginId":id,
                "PluginCapId":$scope.newcap.newcapmaps[i].Id
        };
        AddCapMap_noasync( json, 'POST', '' );
      }

      $scope.okmessage = "The plugin was added."
      $scope.FillPluginTable();
    }).error( function(data,status) {
      if (status>=500) {
        $scope.errtext = "Server said.";
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
  $scope.FillPluginsAvailTable = function(index) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/repoplugins?installable=true"
    }).success( function(data, status, headers, config) {
      $scope.pluginsAvail = data;
      $scope.pluginsAvailFinished = true;
      if( index >= 0 ) {
        // Delete from the meta table
        $scope.pluginsAvailMeta.splice(index,1);
      }
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
  $scope.FillPluginTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/plugins"
    }).success( function(data, status, headers, config) {
      $scope.plugins = data;
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

  $scope.FillPluginTable();

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.dialog = function (id,Name) {
  // --------------------------------------------------------------------
  // Delete Plugin

    $scope.Name = Name;
    $scope.id = id;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        Name: function () {
          return $scope.Name;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.Name + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, Name) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.Name = Name;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.Edit_DeleteModal = function (id,Code) {
  // --------------------------------------------------------------------

    var modalInstance = $modal.open({
      templateUrl: 'DeletePlugincapmap.html',
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
      templateUrl: 'DeletePlugincapmap.html',
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

