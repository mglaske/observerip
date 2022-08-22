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

# Uses
This application / module can be used in a couple different ways:
- By importing it and using it in your own application
- By compiling it and using it as a REST endpoint for the observer data.

## Using as a module
```
import "github.com/mglaske/observerip"

def main() {
  x := observerip.New(5080)  # This is the address it listens on
  go x.ListenAndServe()
}
```

## Compiling and using as a stand-alone app

/*
USAGE:

x := observerip.New(5080)
x.SetPassthrough(true)
x.ListenAndServe()

*/


