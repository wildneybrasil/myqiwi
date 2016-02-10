'use strict';

angular.module('techPayFilters', [])
	.filter('dinheiro',[ '$filter', function(filter) {
    		var currencyFilter = filter('currency');
    		return function(amount, currencySymbol) {
		    if( amount == 0 ) return "";
        		return currencyFilter(amount, currencySymbol);
	      }
	    } ])

  