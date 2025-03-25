import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMap,
  useMapEvent,
} from "react-leaflet";

export interface MapCoords {
  lat: number;
  lng: number;
}

interface MapProps {
  center: MapCoords | undefined;
  setCenter: React.Dispatch<React.SetStateAction<MapCoords | undefined>>;
  disabled?: boolean;
}

export default function MapElem({
  center,
  setCenter,
  disabled = false,
}: MapProps) {
  interface SetLocationProps {
    setCenter: React.Dispatch<React.SetStateAction<MapCoords | undefined>>;
  }

  function SetViewOnClick(props: SetLocationProps) {
    useMapEvent("click", (e) => {
      props.setCenter(e.latlng);
    });

    return null;
  }

  interface LocationProps {
    center: MapCoords;
  }

  function ChangeLocation(props: LocationProps) {
    const map = useMap();
    map.setView(props.center, map.getZoom(), {
      animate: true,
    });

    return (
      <Marker position={props.center}>
        <Popup>{`${props.center.lat}, ${props.center.lng}`}</Popup>
      </Marker>
    );
  }

  return (
    <div>
      {!center || disabled ? (
        <div className="h-[280px] z-0 bg-secondary-gray rounded-lg"></div>
      ) : (
        <MapContainer
          className="h-[280px] z-0 rounded-lg"
          center={center}
          zoom={13}
          scrollWheelZoom={false}
        >
          <TileLayer
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <SetViewOnClick setCenter={setCenter} />
          <ChangeLocation center={center} />
        </MapContainer>
      )}
    </div>
  );
}
