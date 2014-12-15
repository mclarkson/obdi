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
mgrApp.directive("updcapsel", ['$timeout', function ( $timeout ) {
// ------------------------------------------------------------------------
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
      scope.$watch("dccaps", function (value) {
          $timeout( function() {
              element.selectpicker('val',element.val());
              element.selectpicker('refresh');
              //$('.selectpicker').selectpicker('refresh');
          });
      });
    }
}]);

// ------------------------------------------------------------------------
mgrApp.controller("dcCtrl", function ($log, $modal, $scope, $http, baseUrl) {
// ------------------------------------------------------------------------

  $scope.adddc = false;
  $scope.editdc = false;
  $scope.managecaps = false;
  $scope.dccapmaps = [];
  $scope.dccaps = [];
  $scope.newcap = {};
  $scope.newcap.newcapmaps = [];
  $scope.newcap.selected = {};
  $scope.dc = {};
  $scope.dc.SysName = "";
  $scope.delcaps = {};
  $scope.delcaps.ids = [];

  $scope.adddctabs = [
    { title:'Data Centre Details', content:'frag/adddc-detailstab.html'},
    { title:'Capabilities', content:'frag/adddc-capstab.html'}
    ];

  $scope.editdctabs = [
    { title:'Data Centre Details', content:'frag/editdc-detailstab.html'},
    { title:'Capabilities', content:'frag/editdc-capstab.html'}
    ];

  // ----------------------------------------------------------------------
  $scope.$watch('dc.SysName', function() {
  // ----------------------------------------------------------------------
    $scope.dc.SysName = $scope.dc.SysName.replace(/\s+/g,'');
  });

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
  $scope.AddDC = function(tf) {
  // ----------------------------------------------------------------------
    $scope.adddc = tf;
    $scope.newcap = {};
    $scope.newcap.newcapmaps = [];
    $scope.newcap.selected = {};
    $scope.dc = {};
    $scope.dc.SysName = "";
    $scope.FillDcCapsTable();
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.EditAddCapMap = function() {
  // ----------------------------------------------------------------------
    $scope.delcaps = {};
    $scope.delcaps.ids = [];

    if( typeof $scope.newcap.selected !== 'undefined' &&
        $scope.newcap.selected ) {

      // Search for the item being added

      var found = $.grep( $scope.dccapmaps,
        function(e){ return e.DcCapId == $scope.newcap.selected.Id; });

      // Don't add it if it's there
      if( found.length == 0 ) {
        $scope.newcap.newcapmaps.push($scope.newcap.selected);
        table_item = {
                      "DcCapCode":$scope.newcap.selected.Code,
                      "DcCapDesc":$scope.newcap.selected.Desc,
                      "DcCapId":$scope.newcap.selected.Id,
                      "DcId":$scope.editdc
        };

        $scope.dccapmaps.push( table_item );
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
  $scope.ManageCaps = function(tf) {
  // ----------------------------------------------------------------------
    $scope.managecaps = tf;
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.EditDC = function(id) {
  // ----------------------------------------------------------------------
    $scope.editdc = id;
    $scope.delcaps = {};
    $scope.delcaps.ids = [];

    if (id) {
      $scope.dc = $.grep($scope.dcs, function(e){ return e.Id == id; })[0];
      //$log.info('$scope.dc.Login: ' + $scope.dc.Login);
      $scope.FillDcCapsTable(); // For Capabilities
      $scope.FillDcCapsMapsTable( id ); // For capability maps
      $scope.newcap.newcapmaps = [];
      $scope.newcap.selected = {};
    }

    $scope.FillDCTable();
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dcs/" + id
    }).success( function(data, status, headers, config) {
      $scope.mainokmessage = "The dc was deleted."
      $scope.FillDCTable();
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
    $scope.dccapmaps = $.grep( $scope.dccapmaps,
      function(e){ return e.Id != id; });

    // Delete from the newcap table in case it was just added
    $scope.newcap.newcapmaps = $.grep( $scope.newcap.newcapmaps,
      function(e){ return e.DcCapId != id; });

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
      data: $scope.dc,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dcs/" + $scope.dc.Id
    }).success( function(data, status, headers, config) {

      // Cycle through each entry in $scope.newcap.newcapmap array
      // and REST ADD to /envcapmaps

      for( var i = 0; i < $scope.newcap.newcapmaps.length; i++ ) {
        json = {
                "DcId":$scope.editdc,
                "DcCapId":$scope.newcap.newcapmaps[i].Id
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
      $scope.FillDCTable();

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

    // Cycle through each entry in $scope.newcap.newcapmap array
    // and REST add to /envcapmaps

    /* */
  }

  // ----------------------------------------------------------------------
  function AddCapMap_noasync( json_obj, type, id ) {
  // ----------------------------------------------------------------------

    json = JSON.stringify( json_obj )
    jQuery.ajax({
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccapmaps" + id,
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
      data: $scope.dc,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dcs"
    }).success( function(data, status, headers, config) {

      id = data.Id;

      // Cycle through each entry in $scope.newcap.newcapmap array
      // and REST add to /envcapmaps

      for( var i = 0; i < $scope.newcap.newcapmaps.length; i++ ) {
        json = {
                "DcId":id,
                "DcCapId":$scope.newcap.newcapmaps[i].Id
        };
        AddCapMap_noasync( json, 'POST', '' );
      }

      $scope.okmessage = "The dc was added."
      $scope.FillDCTable();
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
  $scope.FillDcCapsTable = function() {
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

  // ----------------------------------------------------------------------
  $scope.FillDcCapsMapsTable = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/dccapmaps?dc_id=" + id
    }).success( function(data, status, headers, config) {
      $scope.dccapmaps = data;
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
  $scope.FillDCTable = function() {
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

  $scope.FillDCTable();
  $scope.FillDcCapsTable(); // For Capabilities dropdown

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.dialog = function (id,SysName) {
  // --------------------------------------------------------------------
  // Delete DC

    $scope.SysName = SysName;
    $scope.id = id;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        SysName: function () {
          return $scope.SysName;
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
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, SysName) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.SysName = SysName;

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
      templateUrl: 'DeleteDccapmap.html',
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
      templateUrl: 'DeleteDccapmap.html',
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

