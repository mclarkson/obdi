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

// Set a global so we know what interface we're in (admin or user)
window.interface = "run";

// ------------------------------------------------------------------------
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("runCtrl", function ($scope) {

	$scope.VERSION = window.VERSION;

  //$scope.view = "plugins/dashboard/html/view.html"
  $scope.view = ""

  $scope.setView = function( v ) {
    $scope.view = v;
  }

  $scope.whichView = function() {
    return $scope.view;
  }

  // ------------------------------------------------------------------------
  // Resizing of display
  // ------------------------------------------------------------------------

  topOffset = 66;
	width = (window.innerWidth > 0) ?
		window.innerWidth : screen.width;
	if (width < 768) {
		topOffset = 69;
	}
  height = (window.innerHeight > 0) ?
    window.innerHeight : screen.height;
  height = height - topOffset;
  if (height < 1) height = 1;
  if (height > topOffset) {
    $("#page-wrapper").css("min-height", (height) + "px");
  }

  // Loads the correct sidebar on window load,
  // collapses the sidebar on window resize.
  // Sets the min-height of #page-wrapper to window size
  $(function() {
    $(window).bind("load resize", function() {
      topOffset = 66;
      width = (this.window.innerWidth > 0) ?
        this.window.innerWidth : this.screen.width;
      if (width < 768) {
        //$('div.navbar-collapse').addClass('collapse')
        topOffset = 69; // 2-row-menu
      } else {
        //$('div.navbar-collapse').removeClass('collapse')
      }

      height = (this.window.innerHeight > 0) ?
        this.window.innerHeight : this.screen.height;
      height = height - topOffset;
      if (height < 1) height = 1;
      if (height > topOffset) {
        $("#page-wrapper").css("min-height", (height) + "px");
      }
    })
  })

});
