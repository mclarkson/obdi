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
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("saltkeymgrCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.environments = [];
  $scope.keylist = {};
  $scope.keyfilter = "";
  $scope.env = {};
  $scope.status = {};  // For env chooser button
  $scope.grains = [];
  $scope.forminvalid = true; // For grains setting (dc,env,ver)

  // Alerting
  $scope.message = "";
  $scope.mainmessage = "";
  $scope.okmessage = "";
  $scope.login.error = false;

  // Button enabling/disabling and content showing/hiding vars
  $scope.envchosen = {};
  $scope.envchosen.shown = false;
  $scope.envsetting = {};
  $scope.envsetting.shown = false;
  $scope.envsetting.dc = "";
  $scope.envsetting.env = "";
  $scope.envsetting.version = "";
  $scope.envsetting.numupdated = 0;
  $scope.listbtnpressed = false;
  $scope.btnenvlistdisabled = false;
  $scope.showkeybtnblockhidden = false;
  $scope.btnshowkeysdisabled = true;
  $scope.keylist_ready = false;
  $scope.keylist_empty = true;
  $scope.keylist_accept_empty = true;
  $scope.keylist_reject_empty = true;
  $scope.keylist_unaccepted_empty = true;

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.mainmessage = "";
    $scope.okmessage = "";
    $scope.login.error = false;
    $scope.error = false;
  }

  // ----------------------------------------------------------------------
  $scope.copyToController = function( isit ) {
  // ----------------------------------------------------------------------
      $scope.forminvalid = isit;
  }

  // ----------------------------------------------------------------------
  $scope.envChoice = function( envobj, $event ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $event.preventDefault();
    $event.stopPropagation();
    $scope.status.isopen = !$scope.status.isopen;
    $scope.envchosen.shown = true;
    $scope.btnshowkeysdisabled = false;
    $scope.btnenvlistdisabled = true;
    $scope.env = envobj;
  };

  // KEY MANAGEMENT

  // ----------------------------------------------------------------------
  $scope.Accept = function( name ) {
  // ----------------------------------------------------------------------
    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys?hostname=" + name
           + "&type=accept"
           + "&env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
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
        clearMessages();
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
  $scope.Reject = function( name ) {
  // ----------------------------------------------------------------------
    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys?hostname=" + name
           + "&type=reject"
           + "&env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
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
        clearMessages();
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
  $scope.Delete = function( name ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys/" + name
           + "?env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
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
        clearMessages();
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
  $scope.GetKeyOutputLine = function( id ) {
  // ----------------------------------------------------------------------
  // Actually don't bother with the result, just refresh the list
  // since we know Salt has finished now.

    $scope.FillKeyListTable();
  }

  // ----------------------------------------------------------------------
  $scope.GetKeyListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {

      $scope.showkeybtnblockhidden = true;

      var keylist = $.parseJSON(data[0].Text);

      $scope.keylist = keylist;

      if( $scope.keylist.minions.length == 0 ) {
        $scope.keylist_accept_empty = true;
      } else {
        $scope.keylist_accept_empty = false;
      }
      if( $scope.keylist.minions_pre.length == 0 ) {
        $scope.keylist_unaccepted_empty = true;
      } else {
        $scope.keylist_unaccepted_empty = false;
      }
      if( $scope.keylist.minions_rejected.length == 0 ) {
        $scope.keylist_reject_empty = true;
      } else {
        $scope.keylist_reject_empty = false;
      }

      $scope.keylist_ready = true;
      $scope.keylist_empty = false;

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
        clearMessages();
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

  // DC, ENV & VERSION

  // ----------------------------------------------------------------------
  $scope.GetServerSettingOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    //$scope.okmessage = "Server configuration was updated successfully.";
    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {

      try {
        var result = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( result.length == 0  ||
          typeof( result[$scope.envsetting.saltid] ) == undefined ) {
        $scope.message = "The configuration did not complete.";
        $scope.message_jobid = id;
      }

      $scope.envsetting.numupdated += 1;
      if( $scope.envsetting.numupdated == 3 ) {
        $scope.okmessage = "Server configuration was updated successfully.";
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
  };

  // ----------------------------------------------------------------------
  $scope.ApplyGrain = function( saltid, grain, data ) {
  // ----------------------------------------------------------------------
  // Send { Grain:"version",Text:"0.1.2" }

    var config = {};
    config.Grain = grain;
    config.Text = data;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString(),
      data: config
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish( data.JobId, 50, 0,
        $scope.GetServerSettingOutputLine );
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
        clearMessages();
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
  $scope.ApplySettings = function() {
  // ----------------------------------------------------------------------

    clearMessages();

    $scope.envsetting.numupdated = 0;

    if( $scope.envsetting.dc ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "dc",
          $scope.envsetting.dc );
    } else {
      $scope.envsetting.numupdated += 1;
    }

    if( $scope.envsetting.env ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "env",
          $scope.envsetting.env );
    } else {
      $scope.envsetting.numupdated += 1;
    }

    if( $scope.envsetting.version ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "version",
          $scope.envsetting.version );
    } else {
      $scope.envsetting.numupdated += 1;
    }
  }

  // ENV MANAGEMENT

  // ----------------------------------------------------------------------
  $scope.EnvConfig = function( name ) {
  // ----------------------------------------------------------------------

    clearMessages();
    $scope.envchosen.shown = false;
    $scope.envsetting.shown = true;
    $scope.envsetting.gotgrains = false;
    $scope.envsetting.saltid = name;

    $scope.FillGrainsTable( name );
  }

  // ----------------------------------------------------------------------
  $scope.GetGrainsListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {

      var grains = [];

      // Extract data into array
      //
      try {
        grains = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      // The first (and only) item in the array is the hostname
      // Save it to use later.

      var names = [];

      for(var key in grains){
        names.push(key);
      }

      // Save the key/value under the hostname in a new array.

      i = 0;
      for(var key in grains[names[0]]){
        $scope.grains[i] = {};
        $scope.grains[i].key = key;
        $scope.grains[i].value = grains[names[0]][key];
        i = i + 1;
      }

      // Sort the array by key

      $scope.grains.sort(function(a, b){
        if(a.key < b.key) return -1;
        if(a.key > b.key) return 1;
        return 0;
      });

      // Let the view know

      $scope.envsetting.gotgrains = true;

      // DC

      var item = $.grep($scope.grains,
        function(e){ return e.key == 'dc'; });

      if( item.length > 0 ) {
        $scope.envsetting.dc = item[0].value;
      }

      // ENV

      var item = $.grep($scope.grains,
        function(e){ return e.key == 'env'; });

      if( item.length > 0 ) {
        $scope.envsetting.env = item[0].value;
      }

      // VERSION

      var item = $.grep($scope.grains,
        function(e){ return e.key == 'version'; });

      if( item.length > 0 ) {
        $scope.envsetting.version = item[0].value;
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
  };

  // ----------------------------------------------------------------------
  $scope.FillGrainsTable = function( saltid ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,50,0,$scope.GetGrainsListOutputLine);
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
        clearMessages();
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
  $scope.GoBack = function( ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.envchosen.shown = true;
    $scope.envsetting.shown = false;
    $scope.envsetting.dc = "";
    $scope.envsetting.env = "";
    $scope.envsetting.version = "";
    $rootScope.$broadcast( "setsearchtext", $scope.hostfilter );
  }

  // ----------------------------------------------------------------------
  $scope.Restart = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.envchosen.shown = false;
    $scope.listbtnpressed = false;
    $scope.btnenvlistdisabled = false;
    $scope.showkeybtnblockhidden = false;
    $scope.btnshowkeysdisabled = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = true;
  };

  // ----------------------------------------------------------------------
  $scope.KeyList = function() {
  // ----------------------------------------------------------------------
    $scope.btnshowkeysdisabled = true;
    $scope.listbtnpressed = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = false;

    $scope.FillKeyListTable();
  };

  // ----------------------------------------------------------------------
  $scope.showOutputlines = function( id ) {
  // ----------------------------------------------------------------------
  // Redirect the user to the Jobs->Outputlines plugin

    $rootScope.outputlines_plugin = {};
    $rootScope.outputlines_plugin.id = id;
    $scope.setView( "plugins/jobs/html/outputlines.html" );
  }

  // ----------------------------------------------------------------------
  $scope.PollForJobFinish = function( id,delay,count,func ) {
  // ----------------------------------------------------------------------
      $timeout( function() {
        $http({
          method: 'GET',
          url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
               + "/jobs?job_id=" + id
							 + '&time='+new Date().getTime().toString()
        }).success( function(data, status, headers, config) {
          job = data[0];
          if(job.Status == 0 || job.Status == 1 || job.Status == 4) {
            if( count > 40 ) {
              clearMessages();
              $scope.message = "Job took too long. check job ID " +
                               + id + ", then try again.";
              $scope.message_jobid = job['Id'];
            } else {
              // Then retry: capped exponential backoff
              delay = delay < 600 ? delay * 2 : 1000;
              count = count + 1;
              $scope.PollForJobFinish(id,delay,count,func);
            }
          } else if(job.Status == 5) { // Job was successfully completed
            func( id );
          } else { // Some error
            clearMessages();
            $scope.message = "Server said: " + job['StatusReason'];
            $scope.message_jobid = job['Id'];
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
            clearMessages();
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
      }, delay );
  };

  // ----------------------------------------------------------------------
  $scope.FillKeyListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys"
           + "?env_id=" + $scope.env.Id
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyListOutputLine);
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
        clearMessages();
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
  $scope.FillEnvironmentsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs?writeable=1"
           + '&time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.environments = data;
      if( data.length == 0 ) {
        $scope.serverlist_empty = true;
        $scope.btnenvlistdisabled = true;
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
        clearMessages();
        $scope.mainmessage = "Server said: " + data['Error'];
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

  $scope.FillEnvironmentsTable();

  // --------------------------------------------------------------------
  $scope.dialog = function (servername) {
  // --------------------------------------------------------------------

    $scope.servername = servername;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the servername variable is passed to the ModalInstanceCtrl
        servername: function () {
          return $scope.servername;
        }
      }
    });

    modalInstance.result.then(function (servername) {
      $log.info('Will delete: ' + $scope.servername + '(' + $scope.servername + ')' );
      $scope.Delete($scope.servername);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, servername) {
  // --------------------------------------------------------------------

    // So the template can access 'servername' in this new scope
    $scope.servername = servername;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

});
