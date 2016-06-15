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
mgrApp.directive('onLastRepeat', function () {
// ------------------------------------------------------------------------
  return function(scope, element, attrs) {
      if (scope.$last) setTimeout(function(){
        scope.$emit('onRepeatLast', element, attrs);
    }, 1000);
  };
});

// ------------------------------------------------------------------------
mgrApp.directive('ngEnter', function () {
// ------------------------------------------------------------------------
// This directive allows us to pass a function in on an enter key.

  return function (scope, element, attrs) {
    element.bind("keydown keypress", function (event) {
      if(event.which === 13) {
        scope.$apply(function (){
            scope.$eval(attrs.ngEnter);
        });

        event.preventDefault();
      }
    });
  };
});

// ------------------------------------------------------------------------
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("sidebarCtrl", function ($scope,$http,baseUrl,$timeout,
  $rootScope) {

  $scope.plugins = [];
  $scope.sidebarItems = [];
  $scope.sidebarshow = false;
  $scope.text = "";

  $scope.sidebarClick = function( name,hasview ) {
    if( hasview == 1 ) {
      $scope.setView( "plugins/"+name+"/html/view.html" );
    }
  }

  $scope.search = function( text ) {
    $rootScope.$broadcast( "search", text );
  }

  // ----------------------------------------------------------------------
  $scope.$on( "setsearchtext", function( event, args ) {
  // ----------------------------------------------------------------------
    $scope.text = args;
  });

  // ----------------------------------------------------------------------
  $scope.$on( "searchdisabled", function( event, args ) {
  // ----------------------------------------------------------------------
    $scope.searchdisabled = args;
  });

  // ----------------------------------------------------------------------
  $scope.$on('onRepeatLast', function(scope, element, attrs){
  // ----------------------------------------------------------------------
  // Kicked off by the $last ng-repeat. Doesn't work though hence the
  // 500ms delay in the directive.

    $('#side-menu').metisMenu();
    $scope.sidebarshow = true;
    $scope.$apply();
  });

  // ----------------------------------------------------------------------
  $scope.CreateSidebarArray = function() {
  // ----------------------------------------------------------------------
  /* sidebarItems looks like:
   
      $scope.sidebarItems =
      [
        {
            "name":"dashboard"
            "template":"plugins/dashboard/html/sidebar.html",
            "level2items":[
              {
                "name": "dash-jobs",
                "template":"plugins/dash-jobs/html/sidebar.html"
              },
              {
                "name": "dash-salt",
                "template":"plugins/dash-salt/html/sidebar.html"
              }
            ]
        },
        {
            "name":"jobs"
            "template":"plugins/jobs/html/sidebar.html",
            "level2items":[]
        }
      ];
  */

    var p = $scope.plugins;
    var s = $scope.sidebarItems;

    // Get items without a parent (top level items) first
    for( var i=0; i < p.length; i++ ) {
      if( p[i].Parent.length == 0) {
        if( p[i].Name == "" || p[i].HasView == 0 ) continue;
        a = { name:        p[i].Name,
              template:    "plugins/"+p[i].Name+"/html/sidebar.html",
              hasview:     p[i].HasView,
              level2items: [] };
        s.push(a);
      }
    }
    // Add level 2 items
    for( var i=0; i < p.length; i++ ) {
      if( p[i].Parent.length > 0) {
        if( p[i].Name == "" || p[i].HasView == 0  ) continue;
        // Find the parent
        for( var j=0; j < s.length; j++ ) {
          if( p[i].Parent == s[j].name ) {
            a = { name:     p[i].Name,
                  hasview:  p[i].HasView,
                  template: "plugins/"+p[i].Name+"/html/sidebar.html" };
            s[j].level2items.push(a);
          }
        }
      }
    }
  }

  // ----------------------------------------------------------------------
  $scope.FillPluginsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/plugins"
			     + '?time='+new Date().getTime().toString()
    }).success( function(data, status, headers, config) {
      $scope.plugins = data;
      $scope.CreateSidebarArray();
      /*
      $timeout( function() {
          $('#side-menu').metisMenu();
          $scope.sidebarshow = true;
      }, 1000 );
      */

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

  $scope.FillPluginsTable();
});
