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

mgrApp.controller("saltregexmgrCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.environments = [];
  $scope.regexlist = {};
  $scope.keyfilter = "";
  $scope.mapfilter = "";
  $scope.env = {};
  $scope.status = {};  // For env chooser button

  // Alerting
  $scope.message = "";
  $scope.mainmessage = "";
  $scope.okmessage = "";
  $scope.login.error = false;

  // Button enabling/disabling and content showing/hiding vars
  $scope.envchosen = {};
  $scope.envchosen.shown = false;
  $scope.mapconfig = {};
  $scope.mapconfig.map = [];
  $scope.mapconfig.shown = false;
  $scope.mapconfig.newclass = "";
  $scope.mapconfig.maplist_empty = true;
  $scope.mapconfig.maplist_ready = false;
  $scope.mapconfig.regex_id = 0;
  $scope.mapconfig.regx_name = "";
  $scope.mapconfig.apply_disabled = true;
  $scope.editregex = {};
  $scope.editregex.shown = false;
  $scope.editregex.apply_disabled = false;
  $scope.listbtnpressed = false;
  $scope.btnenvlistdisabled = false;
  $scope.showkeybtnblockhidden = false;
  $scope.btnshowkeysdisabled = true;
  $scope.regexlist_ready = false;
  $scope.regexlist_empty = true;
  $scope.spacing = 20;
  $scope.newregex = {};
  $scope.newregex.Name = ""; // so watch works without error

  $rootScope.$broadcast( "searchdisabled", true );

  // ----------------------------------------------------------------------
  $scope.$watch('newregex.Name', function() {
  // ----------------------------------------------------------------------
    $scope.newregex.Name = $scope.newregex.Name.replace(/\s+/g,'');
  });

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

  // ----------------------------------------------------------------------
  $scope.DeleteRegex = function( name ) {
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
  $scope.copyToController = function( invalid ) {
  // ----------------------------------------------------------------------

    $scope.forminvalid = invalid;

  }

  // ----------------------------------------------------------------------
  $scope.DeleteRegexRest = function( id ) {
  // ----------------------------------------------------------------------

    clearMessages();

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regexes/" + id
           + "?env_id=" + $scope.env.Id,
    }).success( function(data, status, headers, config) {
        $scope.FillRegexListTable();
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
  $scope.ApplyRegex = function() {
  // ----------------------------------------------------------------------

    clearMessages();

    $scope.editregex.apply_disabled = true;

    if( $scope.editregex.newregex == true ) {
        method = 'POST';
    } else {
        method = 'PUT';
    }

    $http({
      method: method,
      data: $scope.newregex,
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regexes"
           + "?env_id=" + $scope.env.Id,
    }).success( function(data, status, headers, config) {
      $scope.okmessage = "Regex configuration was updated successfully.";
      $scope.editregex.apply_disabled = false;
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
  $scope.AddClass = function( classname ) {
  // ----------------------------------------------------------------------
    clearMessages();

    if( !$scope.mapconfig.map ) {
        $scope.mapconfig = {};
        $scope.mapconfig.map = [];
    }

    // Don't allow duplicates
    for( var i=0; i < $scope.mapconfig.map.length; ++i ) {
      if(  $scope.mapconfig.map[i].Class == classname ) {
        // It's in the list so leave
        $scope.message = "Cannot add duplicate class.";
        return;
      }
    }

    // Add it
    //$scope.mapconfig.changed = true;
    //$scope.mapconfig.classeschanged = true;
    $scope.mapconfig.newclass = ""; // reset select box
    $scope.mapconfig.map.push( {"Class":classname} );
    $scope.mapconfig.apply_disabled = false;
  };

  // ----------------------------------------------------------------------
  $scope.ApplyMap = function() {
  // ----------------------------------------------------------------------

    clearMessages();

    // Disable the Apply button
    $scope.mapconfig.apply_disabled = true;

    // Load 'config' with:
    //    Classes     []string
    //    RegexId     int64

    var config = {};
    config.Classes = [];

    // Classes
    map = $scope.mapconfig.map;
    for( var i=0; i < map.length; ++i ) {
      config.Classes.push( map[i].Class);
    }

    // RegexId
    config.RegexId = $scope.mapconfig.regex_id;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regex_sls_maps"
           + "?env_id=" + $scope.env.Id,
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
  $scope.EditRegex = function( index ) {
  // ----------------------------------------------------------------------

    clearMessages();

    $scope.editregex.newregex = false;

    $scope.newregex = $scope.regexlist[index];
    $scope.editregex.title = "Edit settings for '" + $scope.newregex.Name
        + "'";

    $scope.editregex.apply_disabled = false;
    $scope.forminvalid = true;
    $scope.editregex.index = index;

    $scope.envchosen.shown = false;
    $scope.editregex.shown = true;
  }

  // ----------------------------------------------------------------------
  $scope.NewRegex = function() {
  // ----------------------------------------------------------------------
  // Reuse editregex

    clearMessages();

    $scope.editregex.newregex = true;

    $scope.newregex = {};
    $scope.newregex.Name = ""; // so watch works without error
    $scope.editregex.title = "Enter details for the new regular expression:"

    $scope.editregex.apply_disabled = false;
    $scope.forminvalid = true;

    $scope.envchosen.shown = false;
    $scope.editregex.shown = true;
  }

  // ----------------------------------------------------------------------
  $scope.MapConfig = function( regex_id, name ) {
  // ----------------------------------------------------------------------

    clearMessages();
    $scope.envchosen.shown = false;
    $scope.mapconfig.shown = true;
    $scope.mapconfig.regex_id = regex_id;
    $scope.mapconfig.regx_name = name;

    $scope.FillMapsTable( regex_id );
  }

  // ----------------------------------------------------------------------
  $scope.GetDescription = function( name ) {
  // ----------------------------------------------------------------------

    if( typeof(name) == 'undefined' ) return "";

    // Search in formula's
    if( ! name.split(".")[1] ) {
      var desc = $.grep($scope.statedescs,
        function(e){ return (e.FormulaName==name && !e.StateFileName); });

      if( desc.length > 0 ) {
        if( desc[0].Desc.length > 0 ) {
            return desc[0].Desc;
        } else {
          return "No description available.";
        }
      }
    }

    // Search in state files
    desc = $.grep($scope.statedescs,
      function(e){ return e.StateFileName == name.split(".")[1]; });

    if( desc.length > 0 ) {
      if( desc[0].Desc.length > 0 ) {
          return desc[0].Desc;
      } else {
        return "No description available.";
      }
    }

    if( $scope.get_desc_in_progress ) {
        return "Getting description...";
    } else {
        return "No description available.";
    }

  }

  // ----------------------------------------------------------------------
  $scope.FillMapsTable = function( regex_id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regex_sls_maps?regex_id=" + regex_id
           + "&env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {

      // Extract data into array
      try {
        $scope.mapconfig.map = $.parseJSON(data.JsonData);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( $scope.mapconfig.map.length == 0 ) {
        $scope.mapconfig.maplist_empty = true;
      } else {
        $scope.mapconfig.maplist_empty = false;
      }

      for( var i=0; i<$scope.mapconfig.map.length; ++i ) {
          if( $scope.mapconfig.map[i].StateFile == "" ) {
              $scope.mapconfig.map[i].Class = $scope.mapconfig.map[i].Formula;
          } else {
              $scope.mapconfig.map[i].Class = $scope.mapconfig.map[i].Formula +
                  "." + $scope.mapconfig.map[i].StateFile;
          }
      }

      // Let the view know

      $scope.mapconfig.maplist_ready = true;

      // Now load the Class descriptions
      $scope.FillDescriptionTable();

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

    $scope.mapconfig.map = [];
    $scope.mapconfig.shown = false;
    $scope.mapconfig.maplist_empty = true;
    $scope.mapconfig.maplist_ready = false;
    $scope.mapconfig.regex_id = 0;
    $scope.mapconfig.regx_name = "";
    $scope.mapconfig.apply_disabled = true;

    $scope.editregex.shown = false;
    $scope.editregex.index = 0;
    $scope.editregex.apply_disabled = false;
    $scope.newregex = {};
    $scope.newregex.Name = ""; // so watch works without error

    $scope.FillRegexListTable();
    //$rootScope.$broadcast( "setsearchtext", $scope.hostfilter );
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
    $scope.regexlist_ready = false;
    $scope.regexlist_empty = true;
    $scope.spacing = 20;
  };

  // ----------------------------------------------------------------------
  $scope.RegexList = function() {
  // ----------------------------------------------------------------------
    $scope.btnshowkeysdisabled = true;
    $scope.listbtnpressed = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = false;

    $scope.FillRegexListTable();
  };

  // ----------------------------------------------------------------------
  $scope.showOutputlines = function( id ) {
  // ----------------------------------------------------------------------
  // Redirect the user to the Jobs->Outputlines plugin

    $rootScope.outputlines_plugin = {};
    $rootScope.outputlines_plugin.id = id;
    $scope.setView( "plugins/systemjobs/html/outputlines.html" );
  }

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
            if( count > 120 ) {
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
  $scope.GetStatedescOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        $scope.statedescs = $.parseJSON(data[0].Text);
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

      $scope.get_desc_in_progress = false;

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

    $scope.get_desc_in_progress = true;

    $scope.statedescs = [];
    $scope.statedescs_names = [];

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/statedescs"
           + "?env_id=" + $scope.env.Id
           + "&version=0", // Zero version - use main unversioned branch name
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
        $scope.message = "Could not load class list for environment version."
                         + " Server said: " + data['Error'];
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
  $scope.GetDescriptionChooser = function( name ) {
  // ----------------------------------------------------------------------

    if( typeof(name) == 'undefined' ) return "";

    // Search in formula's
    if( ! name.split(".")[1] ) {
      var desc = $.grep($scope.statedescs,
        function(e){ return (e.FormulaName==name && !e.StateFileName); });

      if( desc.length > 0 ) {
        if( desc[0].Desc.length > 0 ) {
            return desc[0].Desc;
        } else {
          return "";
        }
      }
    }

    // Search in state files
    desc = $.grep($scope.statedescs,
      function(e){ return e.StateFileName == name.split(".")[1]; });

    if( desc.length > 0 ) {
      if( desc[0].Desc.length > 0 ) {
          return desc[0].Desc;
      } else {
        return "No description available.";
      }
    }

    if( $scope.get_desc_in_progress ) {
        return "Getting the list of classes...";
    } else {
        return "";
    }

  }

  // ----------------------------------------------------------------------
  $scope.FillRegexListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regexes"
           + "?env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {

      $scope.showkeybtnblockhidden = true;

      var regexlist = $.parseJSON(data.JsonData);

      $scope.regexlist = regexlist;

      if( $scope.regexlist.length == 0 ) {
        $scope.regexlist_empty = true;
      } else {
        $scope.regexlist_empty = false;
      }

      $scope.regexlist_ready = true;
      $scope.spacing = 0;

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
  $scope.DeleteRegex = function (servername,id) {
  // --------------------------------------------------------------------

    $scope.servername = servername;
    $scope.id = id;

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
      $scope.DeleteRegexRest($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.DeleteClassHelper = function (servername) {
  // --------------------------------------------------------------------

     $scope.mapconfig.map = $.grep( $scope.mapconfig.map,
         function( n ) {
           return( n.Class != servername )
         });
      $scope.mapconfig.apply_disabled = false;
  }

  // --------------------------------------------------------------------
  $scope.DeleteClass = function (servername) {
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
      $scope.DeleteClassHelper($scope.servername);
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
