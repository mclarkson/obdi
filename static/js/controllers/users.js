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
mgrApp.controller("userCtrl", function ($log, $uibModal, $scope, $http, baseUrl) {
// ------------------------------------------------------------------------

  $scope.adduser = false;
  $scope.edituser = false;
  $scope.perms = {};
  $scope.user = {};
  $scope.envs = {};

  $scope.user.Login = "";
  $scope.user.Passhash = "";
  $scope.user.passwordv = "";

  $scope.addusertabs = [
    { title:'User Details', content:'frag/adduser-detailstab.html'},
    { title:'Permissions', content:'frag/adduser-permstab.html'}];

  $scope.editusertabs = [
    { title:'User Details', content:'frag/edituser-detailstab.html'},
    { title:'Permissions', content:'frag/edituser-permstab.html'}];

  // ----------------------------------------------------------------------
  $scope.$watch('user.Login', function() {
  // ----------------------------------------------------------------------
    $scope.user.Login = $scope.user.Login.toLowerCase().replace(/\s+/g,'');
  });

  // ----------------------------------------------------------------------
  $scope.FillEnvTable = function( callback ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs"
    }).success( function(data, status, headers, config) {
      $scope.envs = data;
      if( $scope.edituser > 0 ) {
        $scope.SetEnvTableFromPerms( $scope.edituser );
      }
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
  }

  // ----------------------------------------------------------------------
  $scope.adminChangeAllow = function() {
  // ----------------------------------------------------------------------
  // For use with the ng-disabled directive.
  // Will return false unless the current user is the 'admin' user.

      if( $scope.user.Login == "admin" ) {
          return true;
      }
      return false;
  }

  // ----------------------------------------------------------------------
  $scope.editUserTabs = function() {
  // ----------------------------------------------------------------------
  // For use in ng-repeat in the edit user tabset.
  // Strip off the Permissions tab for the admin user.

      tabs = $scope.editusertabs;
      if( $scope.user.Login == "admin" ) {
          return tabs.slice(0,1);
      }
      return tabs;
  }

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.mainmessage = "";
    $scope.mainokmessage = "";
  }

  // ----------------------------------------------------------------------
  $scope.AddUser = function(tf) {
  // ----------------------------------------------------------------------
    $scope.adduser = tf;
    $scope.user = {};
    $scope.FillEnvTable();  // For the Permissions tab
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.SetEnvTableFromPerms = function(id) {
  // ----------------------------------------------------------------------
    
    // Read perms for the user into $scope.perms

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/perms?user_id=" + id
    }).success( function(data, status, headers, config) {
      $scope.perms = data;

      // copy the perms into $scope.envs[N].perm

      for (var i = 0; i < $scope.perms.length; i++) {
        for (var j = 0; j < $scope.envs.length; j++) {
          if ( $scope.envs[j].Id == $scope.perms[i].EnvId ) {
            $scope.envs[j].perm = {};
            $scope.envs[j].perm.Enabled = $scope.perms[i].Enabled;
            $scope.envs[j].perm.Writeable = $scope.perms[i].Writeable;
            $scope.envs[j].perm.Id = $scope.perms[i].Id;
          }
        }
      }
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
  }

  // ----------------------------------------------------------------------
  $scope.EditUser = function(id) {
  // ----------------------------------------------------------------------
    $scope.edituser = id;

    //$scope.FillUserTable();

    if (id) {
      // FillEnvTable will also read the perms table if $scope.edituser is set
      $scope.FillEnvTable();
      // Fill 'user' object for edituser-detailstab.html
      $scope.user = $.grep($scope.users, function(e){ return e.Id == id; })[0];
    } else {
      $scope.FillUserTable();
    }

    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/users/" + id
    }).success( function(data, status, headers, config) {
      $scope.mainokmessage = "The user was deleted."
      $scope.FillUserTable();
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
  $scope.EditApply = function() {
  // ----------------------------------------------------------------------

    var arr = $scope.envs;
    var enabled = false;

    if ( ( typeof $scope.user.Passhash != 'undefined' || 
         typeof $scope.user.passwordv != 'undefined' ) &&
         $scope.user.Passhash != $scope.user.passwordv ) {
        clearMessages();
        $scope.message = "The 'Password' and 'Verify' fields don't match.";
        return;
    } else {
      delete $scope.user.passwordv;
      if ($scope.user.Passhash == '' ) {
        delete $scope.user.Passhash;
      }
      clearMessages();
    }

    $http({
      method: 'PUT',
      data: $scope.user,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/users/" + $scope.user.Id
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "Changes were applied."
      // Write the user Permissions tab
      for (var i = 0; i < arr.length; i++) {
        enabled = false;
        if ( typeof arr[i].perm == 'undefined' ) {
          item = null;
        } else {
          item = arr[i].perm;
        }
        if ( item ) {
          // Perm
          if ( arr[i].perm.Enabled == true ) {
            enabled = true;
          } else {
            enabled = false;
          }
        } else {
          enabled = false;
        }
        if ( enabled == true ) {
          //$log.info( i + ': Enabled item' );
          //$log.info('Enabled: ' + arr[i].perm.Enabled );
          //$log.info('Writeable: ' + arr[i].perm.Writeable );
          json = {
                  "Enabled":true,
                  "Writeable":arr[i].perm.Writeable
          };
          jsonpost = {
                  "UserId":data.Id,
                  "EnvId":arr[i].Id,
                  "Enabled":true,
                  "Writeable":arr[i].perm.Writeable
          };
        } else {
          $log.info( i + ': Disabled item' );
          json = {
                  "Enabled":false,
                  "Writeable":false
          };
          jsonpost = {
                  "UserId":data.Id,
                  "EnvId":arr[i].Id,
                  "Enabled":false,
                  "Writeable":false
          };
        }
        if( item ) {
          ApplyPermissionNoAsync( json, 'PUT', "/"+arr[i].perm.Id );
        }
        if( $scope.error == true ) {
          $scope.error = false;
          ApplyPermissionNoAsync( jsonpost, 'POST', '' );
        }
      }
      $scope.FillUserTable();
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
  ApplyPermissionNoAsync = function( json_obj, type, id ) {
  // ----------------------------------------------------------------------

    json = JSON.stringify( json_obj )
    jQuery.ajax({
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/perms" + id,
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

    var arr = $scope.envs;
    var enabled = false;

    if ($scope.user.Passhash != $scope.user.passwordv) {
      $scope.message = "Password verification failed.";
      return;
    } else {
      delete $scope.user.passwordv;
      $scope.message = "";
    }

    // Write the main user details tab

    $http({
      method: 'POST',
      data: $scope.user,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/users"
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "The user was added."
      // Write the user Permissions tab
      for (var i = 0; i < arr.length; i++) {
        enabled = false;
        if ( typeof arr[i].perm == 'undefined' ) {
          item = null;
        } else {
          item = arr[i].perm;
        }
        if ( item ) {
          // Perm
          if ( arr[i].perm.Enabled == true ) {
            enabled = true;
          } else {
            enabled = false;
          }
        } else {
          enabled = false;
        }
        if ( enabled == true ) {
          $log.info( i + ': Enabled item' );
          //$log.info('Enabled: ' + arr[i].perm.Enabled );
          //$log.info('Writeable: ' + arr[i].perm.Writeable );
          json = {
                  "UserId":data.Id,
                  "EnvId":arr[i].Id,
                  "Enabled":true,
                  "Writeable":arr[i].perm.Writeable
          };
        } else {
          //$log.info( i + ': Disabled item' );
          json = {
                  "UserId":data.Id,
                  "EnvId":arr[i].Id,
                  "Enabled":false,
                  "Writeable":false
          };
        }
        ApplyPermissionNoAsync( json, 'POST', '' )
      }
      $scope.FillUserTable();
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
  $scope.FillPermsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/perms"
    }).success( function(data, status, headers, config) {
      $scope.users = data;
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
  $scope.FillUserTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/users"
    }).success( function(data, status, headers, config) {
      $scope.users = data;
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

  $scope.FillUserTable();

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.dialog = function (id,loginname) {
  // --------------------------------------------------------------------

    $scope.loginname = loginname;
    $scope.id = id;

    var modalInstance = $uibModal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        loginname: function () {
          return $scope.loginname;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.loginname + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $uibModalInstance, loginname) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.loginname = loginname;

    $scope.ok = function () {
      $uibModalInstance.close();
    };

    $scope.cancel = function () {
      $uibModalInstance.dismiss('cancel');
    };
  };

});

