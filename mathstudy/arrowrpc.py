# https://arrow.apache.org/docs/python/flight.html#using-the-flight-client
import pyarrow.flight as flight

client = flight.FlightClient(location="grpc://")
action = flight.Action("get_infos", b"")
optyons = flight.FlightCallOptions(timeout=1)
