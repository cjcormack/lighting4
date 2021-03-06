var module = angular.module('FixtureDirective', ['FixtureService']);

module.directive('fixtureDirective', function() {
	return {
		restrict: 'A',
		templateUrl: '/lighting/style/partials/fixture.html',
		scope: {
			fixtureName: "=fixturename"
		},
		controller: function($scope, FixtureService) {
			$scope.fixture = null;

			FixtureService.getFixture($scope.fixtureName).then(function(fixture) {
				$scope.fixture = fixture;

				if (!$scope.fixture.universe) {
					$scope.fixture.universe = 1;
				}

				FixtureService.getFixtureType(fixture.type).then(function(fixtureType) {
					$scope.fixtureType = fixtureType;
				});
			});
		}
	};
});