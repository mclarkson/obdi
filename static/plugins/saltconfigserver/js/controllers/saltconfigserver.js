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
        scope.$watch("config", function (value) {//I change here
            var val = value || null;            
            if (val)
                element.selectpicker('refresh');
                //$('.selectpicker').selectpicker('refresh');
        });
    };
});

// ------------------------------------------------------------------------
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("saltconfigserverCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.statedescs = [];
  $scope.statedescs_names = [];
  $scope.environments = [];
  $scope.env = {};
  $scope.servernames = [];
  $scope.grains = [];
  $scope.grainskeys = [];
  $scope.checkbox_allnone = false;
  $scope.serverstodo = [];
  $scope.hostfilter = "";
  $scope.filteredItems_saved = [];
  $scope.resultlist = [];
  $scope.resultlist_empty = false;
  $scope.resultlist_ready = false;

  // Pages
  $scope.mainview = true;
  $scope.grainsview = {};
  $scope.grainsview.show = false;
  $scope.grainsview.saltid = "";
  $scope.grainsview.gotgrains = false;
  $scope.fromGrainsCache = false;
  $scope.config = {};
  $scope.configs = {};
  $scope.configview = {};
  $scope.configview.changed = false; // for apply button
  $scope.configview.versionchanged = false; // for send logic
  $scope.configview.classeschanged = false; // for send logic
  $scope.configview.newclass = "";
  $scope.configview.description = "";
  $scope.configview.show = false;
  $scope.configview.gotconfig = false;
  $scope.configview.saltid = "";
  $scope.configview.gotgrains = false;
  $scope.reviewpage = {};
  $scope.reviewpage.enabled = false;
  $scope.applypage = {};
  $scope.applypage.enabled = false;
  $scope.applypage.complete = false;
  $scope.changeversionview = {};
  $scope.changeversionview.show = false;

  // Alerting
  $scope.message = "";
  $scope.mainmessage = "";
  $scope.okmessage = "";
  $scope.login.error = false;

  // Button enabling/disabling and content showing/hiding vars
  $scope.btnenvlistdisabled = false;
  $scope.btnlistserversdisabled = true;
  $scope.btnapplydisabled = true;
  $scope.btnreviewdisabled = true;
  $scope.envchosen = false;
  $scope.listbtnpressed = false;
  $scope.serverlist_ready = false;
  $scope.serverlist_empty = true;
  $scope.status = {};

  // ----------------------------------------------------------------------
  $scope.$on( "search", function( event, args ) {
  // ----------------------------------------------------------------------
    if( $scope.grainsview.show == false ) {
      $scope.hostfilter = args;
      $scope.checkbox_allnone = false;
      for( var i=0; i < $scope.servernames.length; i=i+1 ) {
        $scope.servernames[i].Selected = false;
      }
      ReviewBtnStatus();
    } else {
      $scope.grainfilter = args;
    }
  });

  // ----------------------------------------------------------------------
  $scope.copyToController = function (data) {
  // ----------------------------------------------------------------------
  // I can't get to filteredItems unless I copy the data out first.

    $scope.filteredItems_saved = data;
  }

  // ----------------------------------------------------------------------
  function FilterOut() {
  // ----------------------------------------------------------------------

    // Unselect anything that is filtered out of the list
    for( i=0; i < $scope.servernames.length; ++i ) {
      var servername = $scope.servernames[i].Name;
      var server = $.grep($scope.filteredItems_saved,
        function(e){ return e.Name == servername; });
      if( server.length > 0 ) {
        // It's in the list so leave it alone
        continue;
      }
      // servername is not in filteredItems_saved so unselect it
      $scope.servernames[i].Selected = false;
    }
  }

  // ----------------------------------------------------------------------
  function ReviewBtnStatus() {
  // ----------------------------------------------------------------------

    $scope.btnreviewdisabled = true;

    // Show the Review button if something is selected
    for( i=0; i < $scope.servernames.length; ++i ) {
      if( $scope.servernames[i].Selected == true ) {
        $scope.btnreviewdisabled = false;
        break;
      }
    }
  }

  // ----------------------------------------------------------------------
  $scope.GetDescription = function( name ) {
  // ----------------------------------------------------------------------

    if( typeof(name) == 'undefined' ) return "";

    // Search in formula's
    var desc = $.grep($scope.statedescs,
      function(e){ return e.FormulaName == name; });

    if( desc.length > 0 ) {
      return desc[0].Desc;
    }

    // Search in state files
    desc = $.grep($scope.statedescs,
      function(e){ return e.StateFileName == name.split(".")[1]; });

    if( desc.length > 0 ) {
      return desc[0].Desc;
    }

    return "";

  }

  // ----------------------------------------------------------------------
  $scope.SelectAllNone = function( ) {
  // ----------------------------------------------------------------------

    if( $scope.checkbox_allnone == false ) {
      // Select all
      $scope.checkbox_allnone = true;
      for( i=0; i < $scope.servernames.length; i=i+1 ) {
        if( $scope.servernames[i].Responded == true ) {
          $scope.servernames[i].Selected = true;
        }
      }
    } else {
      // Select none
      $scope.checkbox_allnone = false;
      for( i=0; i < $scope.servernames.length; i=i+1 ) {
        if( $scope.servernames[i].Responded == true ) {
          $scope.servernames[i].Selected = false;
        }
      }
    }
    ReviewBtnStatus();
  }

  // ----------------------------------------------------------------------
  $scope.Selected = function( servername ) {
  // ----------------------------------------------------------------------

    var server = $.grep($scope.servernames,
      function(e){ return e.Name == servername; })[0];

    if( server.length == 0 ) {
      // TODO: RETURN AN ERROR (This error is not possible)
      alert("(ERROR 1013) Not found in list: " + servername);
    }

    server.Selected = server.Selected == true ? false:true;

    ReviewBtnStatus();
  }

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
  $scope.Review = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    FilterOut();
    $scope.btnlistserversdisabled = false;
    $scope.btnenvlistdisabled = true;
    $scope.btnapplydisabled = false;
    $scope.btnreviewdisabled = true;

    $scope.serverlist_ready = false;
    $scope.reviewpage.enabled = true;

    j=0;
    for( i=0; i < $scope.servernames.length; ++i ) {
      if( $scope.servernames[i].Selected == true ) {
        $scope.serverstodo[j] = $scope.servernames[i];
        $scope.FillConfigTableMap($scope.servernames[i].Name);
        ++j;
      }
    }
    
  };

  // ----------------------------------------------------------------------
  $scope.Apply = function() {
  // ----------------------------------------------------------------------
    clearMessages();

    // Disable the numbered buttons
    $scope.btnlistserversdisabled = true;
    $scope.btnenvlistdisabled = true;
    $scope.btnreviewdisabled = true;
    $scope.btnapplydisabled = true;

    // Switch pages
    $scope.reviewpage.enabled = false;
    $scope.applypage.enabled = true;

    // Create a list of saltids to send
    var saltids = [];
    for( var i=0; i<$scope.serverstodo.length; ++i ) {
      saltids.push( $scope.serverstodo[i].Name );
    }

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/salthighstate"
           + "?env_id=" + $scope.env.Id,
      data: saltids
    }).success( function(data, status, headers, config) {
      // salthighstate returns a list of salt jobs
      $scope.PollForJobFinish(data.JobId,50,0,$scope.GetJobsListOutputLine);
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
  $scope.ApplyColour = function() {
  // ----------------------------------------------------------------------
    if( $scope.btnapplydisabled == true ) {
      return "btn-primary";
    } else {
      return "btn-success";
    }
  };

  // ----------------------------------------------------------------------
  $scope.ChangeServerVersion = function( server ) {
  // ----------------------------------------------------------------------
    clearMessages();
    saltid = server.Name;

    var grain = $.grep($scope.servernames,
      function(e){ return e.Name == saltid; })[0];

    if( typeof grain == "undefined" ) {
      // TODO: RETURN AN ERROR
      alert("(ERROR 1016) Not found: " + saltid);
    }

    $scope.mainview = false;
    $scope.changeversionview.grain = grain;
    $scope.changeversionview.versionchanged = false;
    $scope.changeversionview.changed = false;
    $scope.changeversionview.saltid = saltid;
    $scope.changeversionview.show = true;

  };

  // ----------------------------------------------------------------------
  $scope.GetVersionOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        var result = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( result.length == 0  ||
          typeof( result[$scope.changeversionview.saltid] ) == undefined ) {

        $scope.message = "The configuration did not complete.";
        $scope.message_jobid = id;

      }

      $scope.okmessage = "Server configuration was updated successfully.";

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
  $scope.ApplyVersion = function() {
  // ----------------------------------------------------------------------
  // Send { Grain:"version",Text:"0.1.2" }

    clearMessages();

    config = {};
    config.Grain = "version";
    config.Text = $scope.changeversionview.grain.Version;

    // Disable the Apply button
    $scope.changeversionview.changed = false;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id,
      data: config
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,50,0,$scope.GetVersionOutputLine);
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
  $scope.AddClass = function( classname ) {
  // ----------------------------------------------------------------------
    clearMessages();

    if( !$scope.config.Classes ) {
        $scope.config = {};
        $scope.config.Classes = [];
        $scope.config.Environment = $scope.env.SysName;
    }

    // Don't allow duplicates
    for( var i=0; i < $scope.config.Classes.length; ++i ) {
      if(  $scope.config.Classes[i] == classname ) {
        // It's in the list so leave
        $scope.message = "Cannot add duplicate class.";
        return;
      }
    }

    // Add it
    $scope.configview.changed = true;
    $scope.configview.classeschanged = true;
    $scope.config.Classes.push( classname );
  };

  // ----------------------------------------------------------------------
  $scope.ApplyConfig = function() {
  // ----------------------------------------------------------------------
  // Send { Classes:[],Dc:"",Environment:"",Version:"0" }

    clearMessages();

    var config = {};

    if( $scope.configview.classeschanged ) {
      config = $scope.config;
      config.Dc = $scope.env.DcSysName;
    }

    // Disable the Apply button
    $scope.configview.changed = false;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/enc?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id,
      data: config
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "Server configuration was updated successfully.";
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
  $scope.Restart = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.btnlistserversdisabled = true;
    $scope.btnenvlistdisabled = false;
    $scope.btnapplydisabled = true;
    $scope.btnreviewdisabled = true;
    $scope.listbtnpressed = false;
    $scope.serverlist_ready = false;
    $scope.serverlist_empty = false;
    $scope.envchosen = false;
    $scope.reviewpage.enabled = false;
    $scope.applypage.enabled = false;
    $scope.applypage.complete = false;
    $scope.resultlist_empty = false;
    $scope.resultlist_ready = false;
  };

  // ----------------------------------------------------------------------
  $scope.ServerList = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.checkbox_allnone = false;
    $scope.btnlistserversdisabled = true;
    $scope.btnenvlistdisabled = true;
    $scope.btnapplydisabled = true;
    $scope.btnreviewdisabled = true;
    $scope.listbtnpressed = true;
    $scope.serverlist_ready = false;
    $scope.serverlist_empty = false;
    $scope.envchosen = true;
    $scope.reviewpage.enabled = false;

    // Reset the selected field
    for( i=0; i < $scope.servernames.length; ++i ) {
       $scope.servernames[i].Selected = false; 
    }

    $scope.serverstodo = [];
    $scope.FillServerListTable();
  };

  // ----------------------------------------------------------------------
  $scope.envChoice = function( envobj, $event ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $event.preventDefault();
    $event.stopPropagation();
    $scope.status.isopen = !$scope.status.isopen;
    $scope.envchosen = true;
    $scope.btnlistserversdisabled = false;
    $scope.btnenvlistdisabled = true;
    $scope.env = envobj;
  };

  // ----------------------------------------------------------------------
  $scope.showOutputlines = function( id ) {
  // ----------------------------------------------------------------------

    $rootScope.outputlines_plugin = {};
    $rootScope.outputlines_plugin.id = id;
    $scope.setView( "plugins/jobs/html/outputlines.html" );
  }

  // ----------------------------------------------------------------------
  $scope.style = function( n,p ) {
  // ----------------------------------------------------------------------
    ret = "Error";

    switch( n ) {
      case true:
        ret = "";
        break;
      case false:
        ret = "danger";
        break;
    }

    switch( p ) {
      case true:
        ret = "success";
        break;
    }

    return ret;
  }

  // ----------------------------------------------------------------------
  $scope.Stringify = function( item ) {
  // ----------------------------------------------------------------------
    var str;

    switch( typeof item ) {
      case "object":
        str = JSON.stringify( item );
        return str.replace( /,/g, ", " );
        break;
      default:
        return item;
    }

  }

  // ----------------------------------------------------------------------
  $scope.GetGrainsListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
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

      $scope.grainsview.gotgrains = true;

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
  $scope.GetJobsListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        $scope.resultlist = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( $scope.resultlist.length == 0 ) {

        $scope.resultlist_empty = true;
        $scope.resultlist_ready = true;

      } else {

        $scope.resultlist_empty = false;
        $scope.resultlist_ready = true;
        $scope.applypage.complete = true;

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
  $scope.GetServerListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        $scope.servernames = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( $scope.servernames.length == 0 ) {

        $scope.serverlist_empty = true;
        $scope.serverlist_ready = true;
        $scope.btnlistserversdisabled = false;

      } else {

        $scope.serverlist_ready = true;
        $scope.serverlist_empty = false;
        $scope.btnlistserversdisabled = false;

        // Copy servernames array to objects
        /*
        for( i=0; i<servernames.length; ++i ) {
          key = Object.keys(servernames[i])
          $scope.servernames[i] = {
            Name: key[0],
            Selected: false
          };
        }
        */
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
  $scope.PollForJobFinish = function( id,delay,count,func ) {
  // ----------------------------------------------------------------------
      $timeout( function() {
        $http({
          method: 'GET',
          url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
               + "/jobs?job_id=" + id
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
  $scope.GoBack = function( ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.mainview = true;
    $scope.grainsview.gotgrains = false;
    $scope.grainsview.show = false;
    $scope.configview.gotconfig = false;
    $scope.configview.gotgrains = false;
    $scope.configview.show = false;
    $scope.changeversionview.show = false;
    $rootScope.$broadcast( "setsearchtext", $scope.hostfilter );
  }

  // ----------------------------------------------------------------------
  $scope.ViewConfig = function( server ) {
  // ----------------------------------------------------------------------
    clearMessages();
    saltid = server.Name;

    var grain = $.grep($scope.servernames,
      function(e){ return e.Name == saltid; })[0];

    if( typeof grain == "undefined" ) {
      // TODO: RETURN AN ERROR
      alert("(ERROR 1015) Not found: " + saltid);
    }

    $scope.configview.grain = grain;

    $scope.configview.versionchanged = false;
    $scope.configview.classeschanged = false;

    $scope.mainview = false;
    $scope.configview.gotconfig = false;
    $scope.configview.changed = false;

    $scope.configview.saltid = saltid;
    $scope.configview.show = true;

    $scope.FillDescriptionTable( $scope.env.SysName );
    $scope.FillConfigTable( saltid );
  }

  // ----------------------------------------------------------------------
  $scope.ViewDetails = function( server ) {
  // ----------------------------------------------------------------------
    clearMessages();
    saltid = server.Name;

    $scope.mainview = false;
    $scope.grainsview.gotgrains = false;

    $scope.grainsview.saltid = saltid;
    $scope.grainsview.show = true;
    if( server.Responded ) {
      $scope.fromGrainsCache = false;
      $scope.FillGrainsTable( saltid );
    } else {
      $scope.fromGrainsCache = true;
      $scope.FillGrainsCacheTable( saltid );
    }

    $scope.grainfilter = "";
    $rootScope.$broadcast( "setsearchtext", $scope.grainfilter );
  }

  // ----------------------------------------------------------------------
  $scope.FillConfigTableMap = function( saltid ) {
  // ----------------------------------------------------------------------

    var grain = $.grep($scope.servernames,
      function(e){ return e.Name == saltid; })[0];

    if( typeof grain == "undefined" ) {
      // TODO: RETURN AN ERROR
      alert("(ERROR 1011) Not found: " + saltid);
    }

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/enc?salt_id=" + saltid
           + "&env=" + $scope.env.SysName
           + "&version=" + grain.Version
           + "&dc=" + $scope.env.DcSysName
    }).success( function(data, status, headers, config) {
      $scope.configs[saltid] = $.parseJSON(data['EncData']);
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
  $scope.FillConfigTable = function( saltid ) {
  // ----------------------------------------------------------------------

    var grain = $.grep($scope.servernames,
      function(e){ return e.Name == saltid; })[0];

    if( typeof grain == "undefined" ) {
      // TODO: RETURN AN ERROR
      alert("(ERROR 1012) Not found: " + saltid);
    }

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/enc?salt_id=" + saltid
           + "&version=" + grain.Version
           + "&env=" + $scope.env.SysName
           + "&dc=" + $scope.env.DcSysName
    }).success( function(data, status, headers, config) {
      $scope.config = $.parseJSON(data['EncData']);
      $scope.configview.gotconfig = true;
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
  $scope.FillGrainsTable = function( saltid ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id
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
  $scope.FillGrainsCacheTable = function( saltid ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/grainscache?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id
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
  $scope.FillServerListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/servers?env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,50,0,$scope.GetServerListOutputLine);
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
  $scope.GetStatedescOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        $scope.statedescs = $.parseJSON(data[0].Text);
        $scope.statedescs.sort(function(a, b){
          if(a < b) return -1;
          if(a > b) return 1;
          return 0;
        });
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      // Create an array of names for a select box
      for( var i=0; i < $scope.statedescs.length; ++i ) {
        $scope.statedescs_names[i] = $scope.statedescs[i].FormulaName;
        if( $scope.statedescs[i].StateFileName.length > 0 ) {
          $scope.statedescs_names[i] += "." +
            $scope.statedescs[i].StateFileName;
        }
      }
      // And sort the list
      $scope.statedescs_names.sort(function(a, b){
        if(a < b) return -1;
        if(a > b) return 1;
        return 0;
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

  // ----------------------------------------------------------------------
  $scope.FillDescriptionTable = function() {
  // ----------------------------------------------------------------------

    var grain = $.grep($scope.servernames,
      function(e){ return e.Name == saltid; })[0];

    if( typeof grain == "undefined" ) {
      // TODO: RETURN AN ERROR
      alert("(ERROR 1011) Not found: " + saltid);
    }

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/statedescs"
           + "?env_id=" + $scope.env.Id
           + "&version=" + grain.Version,
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,50,0,$scope.GetStatedescOutputLine);
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

  // Modal dialog

  // --------------------------------------------------------------------
  $scope.Delete = function (classname) {
  // --------------------------------------------------------------------
     $scope.config.Classes = $.grep( $scope.config.Classes,
         function( n ) {
           return( n != classname )
         });
      $scope.configview.changed = true;
      $scope.configview.classeschanged = true;
  }

  // --------------------------------------------------------------------
  $scope.dialog = function (classname) {
  // --------------------------------------------------------------------

    $scope.classname = classname;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the classname variable is passed to the ModalInstanceCtrl
        classname: function () {
          return $scope.classname;
        }
      }
    });

    modalInstance.result.then(function (classname) {
      $log.info('Will delete: ' + $scope.classname + '(' + $scope.classname + ')' );
      $scope.Delete($scope.classname);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, classname) {
  // --------------------------------------------------------------------

    // So the template can access 'classname' in this new scope
    $scope.classname = classname;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

});
