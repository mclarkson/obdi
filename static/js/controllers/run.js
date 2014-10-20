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

// Set a global so we know what interface we're in (admin or user)
window.interface = "run";

// ------------------------------------------------------------------------
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("runCtrl", function ($scope) {

  $scope.view = "plugins/dashboard/html/view.html"

  $scope.setView = function( v ) {
    $scope.view = v;
  }

  $scope.whichView = function() {
    return $scope.view;
  }

  // ------------------------------------------------------------------------
  // Resizing of display
  // ------------------------------------------------------------------------

  topOffset = 50;
  height = (window.innerHeight > 0) ?
    window.innerHeight : screen.height;
  height = height - topOffset;
  if (height < 1) height = 1;
  if (height > topOffset) {
    $("#page-wrapper").css("min-height", (height-30) + "px");
  }

  // Loads the correct sidebar on window load,
  // collapses the sidebar on window resize.
  // Sets the min-height of #page-wrapper to window size
  $(function() {
    $(window).bind("load resize", function() {
      topOffset = 50;
      width = (this.window.innerWidth > 0) ?
        this.window.innerWidth : this.screen.width;
      if (width < 768) {
        $('div.navbar-collapse').addClass('collapse')
        topOffset = 100; // 2-row-menu
      } else {
        $('div.navbar-collapse').removeClass('collapse')
      }

      height = (this.window.innerHeight > 0) ?
        this.window.innerHeight : this.screen.height;
      height = height - topOffset;
      if (height < 1) height = 1;
      if (height > topOffset) {
        $("#page-wrapper").css("min-height", (height-30) + "px");
      }
    })
  })

});
