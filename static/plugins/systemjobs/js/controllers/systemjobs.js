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

mgrApp.controller("systemjobsCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.jobs = [];
  $scope.jobfilter = "";

  $rootScope.$broadcast( "searchdisabled", false );

  // ----------------------------------------------------------------------
  $scope.$on( "search", function( event, args ) {
  // ----------------------------------------------------------------------
    $scope.jobfilter = args;
  });

  // ----------------------------------------------------------------------
  $scope.jobrunning = function( status ) {
  // ----------------------------------------------------------------------
    return (status == 4 || status == 0);
  }

  // ----------------------------------------------------------------------
  $scope.showOutputlines = function( id ) {
  // ----------------------------------------------------------------------

    $rootScope.outputlines_plugin = {};
    $rootScope.outputlines_plugin.id = id;
    $scope.setView( "plugins/systemjobs/html/outputlines.html" );
  }

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.login.error = false;
    $scope.error = false;
  }

  // ----------------------------------------------------------------------
  $scope.prettyDate = function( d ) {
  // ----------------------------------------------------------------------
    var date = new Date(d);
    return date.toLocaleString();
  }

  // ----------------------------------------------------------------------
  $scope.duration = function( start,end,status ) {
  // ----------------------------------------------------------------------
    var unit = "seconds";
    var s = new Date( start );
    var e = new Date( end );
    if( status == 4 || status == 1 ) {
      e = new Date();
    }
    duration = Math.round((e-s)/1000);
    if( duration == 1 ) unit = "second";
    if( duration > 60 ) {
      duration = Math.round(duration/60)
      unit = "minutes"
      if( duration == 1 ) unit = "minute";
    }
    if( duration > 60 ) {
      duration = Math.round(duration/6)/10
      unit = "hours"
      if( duration == 1 ) unit = "hour";
    }
    return duration + " " + unit;
  }

  // ----------------------------------------------------------------------
  $scope.style = function( n ) {
  // ----------------------------------------------------------------------
    ret = "Error";

    switch( n ) {
      case 0:
        ret = "Unknown";
        break;
      case 1:
        ret = "active";
        break;
      case 2:
        ret = "warning";
        break;
      case 3:
        ret = "warning";
        break;
      case 4:
        ret = "active";
        break;
      case 5:
        ret = "success";
        break;
      case 6:
        ret = "danger"
        break;
    }

    return ret;
  }

  // ----------------------------------------------------------------------
  $scope.lookupErrorCode = function( n ) {
  // ----------------------------------------------------------------------
    ret = "Error";

    switch( n ) {
      case 0:
        ret = "Unknown";
        break;
      case 1:
        ret = "Not started";
        break;
      case 2:
        ret = "User cancelled";
        break;
      case 3:
        ret = "System cancelled";
        break;
      case 4:
        ret = "In progress";
        break;
      case 5:
        ret = "Finished, OK";
        break;
      case 6:
        ret = "Finished, FAIL"
        break;
    }

    return ret;
  }

  // ----------------------------------------------------------------------
  $scope.FillJobsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/jobs"
    }).success( function(data, status, headers, config) {
      // lookup the error codes and add as ErrText field.
      for( var i=0; i<data.length; i++ ) {
        data[i].ErrText = $scope.lookupErrorCode(data[i].Status);
      }
      $scope.jobs = data;
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
  $scope.StopJob = function (id) {
  // --------------------------------------------------------------------

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/jobs/kill/" + id
    }).success( function(data, status, headers, config) {
      $timeout( $scope.FillJobsTable, 2000 );
      //$timeout( $scope.FillJobsTable(), 4000 );
      //$timeout( $scope.FillJobsTable(), 6000 );
      clearMessages();
      $scope.okmessage = "The job was stopped."
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

  // --------------------------------------------------------------------
  $scope.dialog = function (id,Name) {
  // --------------------------------------------------------------------
  // Stop job

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
        },
        id: function () {
          return $scope.id;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.Name + '(' + $scope.id + ')' );
      $scope.StopJob($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, Name, id) {
  // --------------------------------------------------------------------

    // So the template can access 'loginname' in this new scope
    $scope.Name = Name;
    $scope.id = id;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

  // --------------------------------------------------------------------
  $scope.jobdetails = function (id) {
  // --------------------------------------------------------------------
  // Stop job

    //$scope.jobs = Name;
    //$scope.id = id;
    $scope.job = $.grep($scope.jobs, function(e){ return e.Id == id; })[0];

    var modalInstance = $modal.open({
      templateUrl: 'jobdetails.html',
      controller: $scope.JobDetailsInstCtrl,
      size: 'md',
      resolve: {
        // the loginname variable is passed to the ModalInstanceCtrl
        job: function () {
          return $scope.job;
        }
      }
    });

    modalInstance.result.then(function (id) {
      $log.info('Will delete: ' + $scope.Name + '(' + $scope.id + ')' );
      $scope.StopJob($scope.id);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.JobDetailsInstCtrl = function ($scope, $modalInstance, job) {
  // --------------------------------------------------------------------

    $scope.prettyDate = function( d ) {
      var date = new Date(d);
      return date.toLocaleString();
    }

    $scope.lookupErrorCode = function( n ) {
      ret = "Error";
      switch( n ) {
        case 0:
          ret = "Unknown";
          break;
        case 1:
          ret = "Not started";
          break;
        case 2:
          ret = "User cancelled";
          break;
        case 3:
          ret = "System cancelled";
          break;
        case 4:
          ret = "In progress";
          break;
        case 5:
          ret = "Finished, OK";
          break;
        case 6:
          ret = "Finished, FAIL"
          break;
      }
      return ret;
    }

    $scope.lookupJobType = function( n ) {
      ret = "Error";

      switch( n ) {
        case 0:
          ret = "Unknown";
          break;
        case 1:
          ret = "User";
          break;
        case 2:
          ret = "System";
          break;
        default:
          ret = "Invalid type"
          break;
      }
      return ret;
    }

    // So the template can access 'loginname' in this new scope
    $scope.job = job;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

  $scope.FillJobsTable();
});
