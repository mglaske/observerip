Observer IP Module
==================

Ambient Weather's Observer IP module connects to an ethernet network and
receives information from temperature and other weather modules.  It's primary
purpose is to push this information to Ambient Weather and/or weather
underground.  Previous versions of the firmware would allow you to set an
endpoint of your choosing, but the have since removed that functionality.
Therefore, the only way to gather this information is to intercept, and
redirect the pushes that it sends out to somewhere else and collect them.

There are a couple documents on the internet that explain how to do this with
various routers/firewalls.  Explaining how to do this is outside the scope of
this project, however the basics are, for any traffic on port 80 or 443 
coming from the local IP address associated with your observer IP module, use
NAT and send those requests to the box where you're hosting this application.

The application supports passthrough (although this hasn't been tested) to
weather undergound as well.  If SetPassthrough(true) is set, the app will
take whatever is sent from the ObserverIP module, disceminate that information
so you can poll the app via REST for your sensor data, and also spawn a
goroutine to send that data to the weather underground site.

NOTE: use this for what it's worth, I haven't had much time to put into this,
so it's very minimal.

# REST Endpoints

- /endpoints : returns the last endpoint response
- /info : returns the last station response

# Data Output Excamples
## Info / Station Response
```
{"Passkey":"blah..blah","StationType":"WS-1501-IP","DateUTC":"2022-08-22 16:47:09","Endpoints":[{"Id":1,"Name":"","TempF":78.08,"Humidity":61,"Battery":1},{"Id":2,"Name":"","TempF":83.3,"Humidity":56,"Battery":1},{"Id":3,"Name":"","TempF":70.88,"Humidity":61,"Battery":1},{"Id":4,"Name":"","TempF":74.48,"Humidity":53,"Battery":1},{"Id":5,"Name":"","TempF":-1,"Humidity":-1,"Battery":1},{"Id":6,"Name":"","TempF":-1,"Humidity":-1,"Battery":1},{"Id":7,"Name":"","TempF":-1,"Humidity":-1,"Battery":1},{"Id":8,"Name":"","TempF":-1,"Humidity":-1,"Battery":1}]}
```

## Endpoint Response
```
{"Id":"STATIONID","Password":"","TempF":-9999,"Humidity":-9999,"DewPointF":-9999,"WindChillF":-9999,"WindDirection":-9999,"WindSpeedMph":-9999,"WindGustMph":-9999,"RainInch":0,"DailyRainInch":0,"WeeklyRainInch":0,"MonthlyRainInch":0,"YearlyRainInch":0,"SolarRadiation":-9999,"UV":-9999,"IndoorTempF":-9999,"IndoorHumidity":-9999,"BarometricInch":-9999,"LowBattery":0,"DateUTC":"now","SoftwareType":"WH2602 V4.6.2","Action":"updateraw","RealTime":1,"RealTimeFrequency":5}
```

# USES
This application / module can be used in a couple different ways:
- By importing it and using it in your own application
- By compiling it and using it as a REST endpoint for the observer data.

## Using as a module
```
import "github.com/mglaske/observerip"

def main() {
  x := observerip.New(5080)  # This is the address it listens on
  go x.ListenAndServe()
  defer x.Close()

  // The last station response received
  sr := x.GetStationResponse()

  // The last endpoint response received
  er := x.GetEndpointResponse()

}
```

## Compiling and using as a stand-alone app

See the example directory.
