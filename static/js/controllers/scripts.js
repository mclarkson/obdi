// Deployment Manager - a REST interface and GUI for deploying software
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
mgrApp.directive('fileUpload', function () {
// ------------------------------------------------------------------------
    return {
        scope: true,        //create a new scope
        link: function (scope, el, attrs) {
            el.bind('change', function (event) {
                var files = event.target.files;
                //iterate files since 'multiple' may be specified on the element
                for (var i = 0;i<files.length;i++) {
                    //emit event upward
                    scope.$emit("fileSelected", { file: files[i] });
                }                                       
            });
        }
    };
});

// ------------------------------------------------------------------------
mgrApp.controller("scriptCtrl", function ($log, $modal, $scope, $http,
      baseUrl) {
// ------------------------------------------------------------------------

  $scope.addscript = false;
  $scope.editscript = false;
  $scope.scripts = {};
  $scope.eye_enabled = true;

  $scope.addscripttabs = [
    { title:'Script Details', content:'frag/addscript-detailstab.html'}
    ];

  $scope.editscripttabs = [
    { title:'Script Details', content:'frag/editscript-detailstab.html'}
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
  $scope.AddScript = function(tf) {
  // ----------------------------------------------------------------------
    $scope.addscript = tf;
    $scope.script = {};
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.EditScript = function(id) {
  // ----------------------------------------------------------------------
    $scope.editscript = id;
    $scope.eye_enabled = true;

    if (id) {
      $scope.files = [];
      $scope.script = $.grep($scope.scripts, function(e){ return e.Id == id; })[0];
      //$log.info('$scope.dc.Login: ' + $scope.dc.Login);
    }

    $scope.FillScriptsTable();
    clearMessages();
  }

  // ----------------------------------------------------------------------
  $scope.Apply = function() {
  // ----------------------------------------------------------------------
  // Apply button for adding a new user

    clearMessages();
    $scope.script.Source = $scope.file_b64;

    $http({
      method: 'POST',
      data: $scope.script,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/scripts"
    }).success( function(data, status, headers, config) {

      id = data.Id;

      $scope.okmessage = "The script was added."
      $scope.FillScriptsTable();

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
  $scope.DownloadFile = function(name) {
  // ----------------------------------------------------------------------
  // Apply button for editing a new user

    clearMessages();

    if( $scope.files.length == 1 ) {
        $scope.script.Source = $scope.file_b64;
    }

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/scripts?id=" + $scope.script.Id
    }).success( function(data, status, headers, config) {

      var link = document.createElement("A");
      link.download = name;
      link.target = "_blank";
      type = data[0].Type;
      if( /text/i.test(type) ) {
        link.href = "data:text/plain;base64,"+data[0].Source;
      } else {
        link.href = "data:application/octet-stream;base64,"+data[0].Source;
      }
      document.body.appendChild(link);
      link.click();

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
  $scope.EditApply = function() {
  // ----------------------------------------------------------------------
  // Apply button for editing a new user

    clearMessages();

    if( $scope.files.length == 1 ) {
        $scope.script.Source = $scope.file_b64;
    }

    $http({
      method: 'PUT',
      data: $scope.script,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/scripts/" + $scope.script.Id
    }).success( function(data, status, headers, config) {

      $scope.okmessage = "Changes were applied."
      $scope.FillScriptsTable();

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
  $scope.FillScriptsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/scripts?nosource=true"
    }).success( function(data, status, headers, config) {
      $scope.scripts = data;
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

  $scope.FillScriptsTable();

  // ----------------------------------------------------------------------
  $scope.Delete = function( id ) {
  // ----------------------------------------------------------------------

    $scope.mainmessage = "";

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/scripts/" + id
    }).success( function(data, status, headers, config) {
      $scope.mainokmessage = "The script was deleted."
      $scope.FillScriptsTable();
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
  // Modal delete dialog
  // --------------------------------------------------------------------

  // --------------------------------------------------------------------
  $scope.dialog = function (id,scriptname) {
  // --------------------------------------------------------------------
  // Delete DC

    $scope.scriptname = scriptname;
    $scope.id = id;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        scriptname: function () {
          return $scope.scriptname;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.scriptname + '(' + $scope.id + ')' );
      $scope.Delete($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, scriptname) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.scriptname = scriptname;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

  // ----------------------------------------------------------------------
  // File upload support
  // ----------------------------------------------------------------------
  //
  // Supports:
  //   <input type="file" file-upload multiple />

  $scope.files = [];
  $scope.file_b64 = "";

  // ------------------------------------------------------------------------
  handleFileSelect = function(evt,args) {
  // ------------------------------------------------------------------------
    var file = args.file.name;
  
    if (file) {
      var reader = new FileReader();

      reader.onload = function(readerEvt) {
          var binaryString = readerEvt.target.result;
          $scope.file_b64 = btoa(binaryString);
          $scope.$apply();
      };

      reader.readAsBinaryString(args.file);
      $scope.eye_enabled = false;
    }
  };
  
  // ------------------------------------------------------------------------
  $scope.$on("fileSelected", function (event, args) {
  // ------------------------------------------------------------------------
  //listen for the file selected event
    $scope.$apply(function () {
      $scope.files = [];
      $scope.files.push(args.file);
      handleFileSelect( event, args );
    });
  });
  
  // ------------------------------------------------------------------------
  $scope.save = function() {
  // ------------------------------------------------------------------------
  //the save method
    json = {
      name: "name",
      file: $scope.file_b64
    };
    $http({
      method: 'POST',
      url: "/Api/PostStuff",
      data: $scope.json
    }).
    success(function (data, status, headers, config) {
      alert("success!");
    }).
    error(function (data, status, headers, config) {
      alert("failed!");
    });
  };

});

