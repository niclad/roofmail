# Roofmail — An SMS companion dedicated to helping you find the perfect relaxation time

## Helpful links
### Weather API
The Government (currently) provides an API that's free to use. [Info here.](https://www.weather.gov/documentation/services-web-api). Using this, it's possible to get forcast and weather data based on geographic coordinates. However, the resolution of this data is only precise down to an area of 2.5km x 2.5km — which is good enough for our use case here.

### Beaufort Scale
The Beaufort Scale is a handy way to determine how unpleasant wind can be. The NWS has a helpful table [here](https://www.weather.gov/pqr/wind), but I'll provide it below as well.

| Beaufort number 	| Description     	| Speed (mph) 	| Notes                                                                              	|
|----------------:	|-----------------	|-------------	|------------------------------------------------------------------------------------	|
|               0 	| Calm            	| [0-1)       	|                                                                                    	|
|               1 	| Light Air       	| [1-4)       	|                                                                                    	|
|               2 	| Light Breeze    	| [4-8)       	|                                                                                    	|
|               3 	| Gentle Breeze   	| [8-13)      	|                                                                                    	|
|               4 	| Moderate Breeze 	| [13-19)     	|                                                                                    	|
|               5 	| Fresh Breeze    	| [19-25)     	| This is where I'd consider the upper limit of "tolerable" for a windy day outside. Below 50°F, wind chill can feel near freezing 	|
|               6 	| Strong Breeze   	| [25-32)     	|                                                                                    	|
|               7 	| Near Gale       	| [32-39)     	|                                                                                    	|
|               8 	| Gale            	| [39-47)     	|                                                                                    	|
|               9 	| Strong Gale     	| [47-55)     	| Structural damage starts here.                                                     	|
|              10 	| Whole Gale      	| [55-64)     	|                                                                                    	|
|              11 	| Storm Force     	| [64-75]     	| Large trees may be uprooted.                                                       	|
|              12 	| Hurricane Force 	| > 75        	| Severe damage can occur to structures... Roofs blown off, windows broken, etc.     	|

