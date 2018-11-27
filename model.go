package onvif

// Device contains data of ONVIF camera
type Device struct {
	ID       string
	Name     string
	XAddr    string
	User     string
	Password string
}

// DeviceInformation contains information of ONVIF camera
type DeviceInformation struct {
	FirmwareVersion string
	HardwareID      string
	Manufacturer    string
	Model           string
	SerialNumber    string
}

// NetworkInterfaces contains information of ONVIF camera
type NetworkInterfaces struct {
}

// NetworkCapabilities contains networking capabilities of ONVIF camera
type NetworkCapabilities struct {
	DynDNS     bool
	IPFilter   bool
	IPVersion6 bool
	ZeroConfig bool
}

// MediaCapabilities contains media capabilities of ONVIF camera
type MediaCapabilities struct {
	XAddr     string
}

// DeviceCapabilities contains capabilities of an ONVIF camera
type DeviceCapabilities struct {
	Network   NetworkCapabilities
	Media	  MediaCapabilities
	Events    map[string]bool
	Streaming map[string]bool
}

// HostnameInformation contains hostname info of an ONVIF camera
type HostnameInformation struct {
	Name     string
	FromDHCP bool
}

// MediaBounds contains resolution of a video media
type MediaBounds struct {
	Height int
	Width  int
}

// MediaSourceConfig contains configuration of a media source
type MediaSourceConfig struct {
	Name        string
	Token       string
	SourceToken string
	Bounds      MediaBounds
}

// VideoRateControl contains rate control of a video
type VideoRateControl struct {
	BitrateLimit     int
	EncodingInterval int
	FrameRateLimit   int
}

// VideoEncoderConfig contains configuration of a video encoder
type VideoEncoderConfig struct {
	Name           string
	Token          string
	Encoding       string
	Quality        int
	RateControl    VideoRateControl
	Resolution     MediaBounds
	SessionTimeout string
}

// AudioEncoderConfig contains configuration of an audio encoder
type AudioEncoderConfig struct {
	Name           string
	Token          string
	Encoding       string
	Bitrate        int
	SampleRate     int
	SessionTimeout string
}

// PTZConfig contains configuration of a PTZ control in camera
type PTZConfig struct {
	Name      string
	Token     string
	NodeToken string
}

// MediaProfile contains media profile of an ONVIF camera
type MediaProfile struct {
	Name               string
	Token              string
	VideoSourceConfig  MediaSourceConfig
	VideoEncoderConfig VideoEncoderConfig
	AudioSourceConfig  MediaSourceConfig
	AudioEncoderConfig AudioEncoderConfig
	PTZConfig          PTZConfig
}

// MediaURI contains streaming URI of an ONVIF camera
type MediaURI struct {
	URI                 string
	Timeout             string
	InvalidAfterConnect bool
	InvalidAfterReboot  bool
}

// System Date And Time
type Time struct {
	Hour int
	Minute int
	Second int
}

type Date struct {
	Year int
	Month int
	Day int
}

type UTCDateTime struct {
	Time
	Date
}
type TimeZone struct {
	TZ string
}

type SystemDateAndTime struct {
	DateTimeType string
	DaylightSavings bool
	TimeZone
	UTCDateTime
}

//NTP information struct
type NetworkHost struct {
	Type string
	IPv4Address string
	IPv6Address string
	DNSname string
}


type NTPInformation struct {
	FromDHCP bool
	NTPNetworkHost NetworkHost // NetworkHost of NTPFromDHCP if FromDHCP is true, else of NTPManual
}

// DNS Information struct
type IPAddress struct {
	Type string
	IPv4Address string
	IPv6Address string
}


type DNSInformation struct {
	FromDHCP bool
	SearchDomain string
	DNSAddress IPAddress // IPAddress of DNSFromDHCP if FromDHCP is true, else of DNSManual
}
// DynamicDNSInformation
type DynamicDNSInformation struct {
	Type string // 'NoUpdate', 'ClientUpdates', 'ServerUpdates'
	Name string
	TTL string
}

// NetWork Interface
type NetworkInterfaceInfo struct {
	Name string
	HwAddress string
	MTU int
}

type NetworkInterfaceConnectionSetting struct {
	AutoNegotiation bool
	Speed int
	Duplex string // "Full" || "Half"
}

type NetworkInterfaceLink struct {
	AdminSettings NetworkInterfaceConnectionSetting
	OperSettings NetworkInterfaceConnectionSetting
	InterfaceType string
}

type PrefixedIPAdress struct {
	Address string
	PrefixLength int
}

type IPv4Configuration struct {
	Manual PrefixedIPAdress
	LinkLocal PrefixedIPAdress
	FromDHCP PrefixedIPAdress
	DHCP bool
}

type IPv4NetworkInterface struct {
	Enabled bool
	Config IPv4Configuration
}

type IPv6Configuration struct {
	AcceptRouterAdvert bool
	DHCP string // 'Auto' 'Stateful' 'Stateless' 'Off'
	Manual PrefixedIPAdress
	LinkLocal PrefixedIPAdress
	FromDHCP PrefixedIPAdress
	FromRA PrefixedIPAdress
}

type IPv6NetworkInterface struct {
	Enable bool
	Config IPv6Configuration
}

type NetworkInterface struct {
	Token string
	Enabled bool
	Info NetworkInterfaceInfo
	Link NetworkInterfaceLink
	IPv4 IPv4NetworkInterface
	IPv6 IPv6NetworkInterface
}

// NetWork Protocols
type NetworkProtocol struct {
	Name string // 'HTTP' 'HTTPS' 'RTSP'
	Enabled bool
	Port int
}

// Scope
type Scope struct {
	ScopeDes string // 'Fixed' 'Configurabale'
	ScopeItem string
}

// Network GateWay
type NetworkGateway struct {
	IPv4Address string
	IPv6Address string
}

// User
type User struct {
	Username string
	Password string
	UserLevel string // 'Administrator', 'Operator', 'User', 'Anonymous', 'Extended'
}

// RelayOutput
type RelayOutputSettings struct {
	Mode string // 'Bistable', 'Monostable'
	DelayTime string
	IdleState string // 'closed', 'open'
}

type RelayOutput struct {
	Token string
	Properties RelayOutputSettings
}

//NetworkZeroConfiguration
type NetworkZeroConfiguration struct {
	InterfaceToken string
	Enabled bool
	Addresses []string
}

// Service Device
type DeviceSecurityCapabilitiesService struct {
	RemoteUserHandling bool
	Dot1X bool
	AccesssPolicyConfig bool
	OnboardKeyGeneration bool
	HttpDigest bool
	X509Token bool
	DefaultAccessPolicy bool
	RELToken bool
	KerberosToken bool
	TLS12 bool
	TLS11 bool
	TLS10 bool
	UsernameToken bool
	SAMLToken bool
}

type DeviceSystemCapabilitiesService struct {
	DiscoveryBye bool
	DiscoveryResolve bool
	FirmwareUpgrade bool
	SystemLogging bool
	SystemBackup bool
	RemoteDiscovery bool
}

type DeviceNetworkCapabilitiesService struct {
	NTP int
	DynDNS bool
	IPVersion6 bool
	ZeroConfiguration bool
	IPFilter bool
}

type DeviceCapabilitiesService struct {
	Security DeviceSecurityCapabilitiesService
	System DeviceSystemCapabilitiesService
	Network DeviceNetworkCapabilitiesService
}

// Service Media
type MediaProfileCapabilitiesService struct {
	MaximumNumberOfProfiles int
}
type MediaStreamingCapabilitiesService struct {
	RTP_RTSP_TCP bool
	RTP_TCP bool
	RTPMulticast bool
	NoRTSPStreaming bool
	NonAggregateControl bool
}

type MediaCapabilitiesService struct {
	OSD bool
	VideoSourceMode bool
	Rotation bool
	SnapshotUri bool
	ProfileCapabilities MediaProfileCapabilitiesService
	StreamingCapabilities MediaStreamingCapabilitiesService
}

// Service Events
type EventsCapabilitiesService struct {
	MaxNotificationProducers int
	WSPausableSubscriptionManagerInterfaceSupport bool
	WSPullPointSupport bool
	WSSubscriptionPolicySupport bool
	PersistentNotificationStorage bool
	MaxPullPoints int
}

// Service Imaging
type ImagingCapabilitiesService struct {
	ImageStabilization bool
}

// Service PTZ
type PTZCapabilitiesService struct {
	GetCompatibleConfigurations bool
	Reverse bool
	EFlip bool
}

type OnvifVersion struct {
	Major int
	Minor int
}

type CapabilitiesService struct {
	Name string // PTZ, media, device, imaging, events
	Capabilities interface{}
}

type Service struct {
	Namespace string
	XAddr string
	Capabilities CapabilitiesService
	Version OnvifVersion
}


// VideoEncoderConfigurationOptions
type IntRange struct {
	Min int
	Max int
}


type H264Options struct {
	ResolutionsAvailable MediaBounds
	GovLengthRange IntRange
	FrameRateRange IntRange
	EncodingIntervalRange IntRange
	H264ProfilesSupported []string // 'Baseline', 'Main', 'Extended', 'High'
}

type VideoEncoderConfigurationOptions struct {
	QualityRange IntRange
	H264 H264Options
}

// GuaranteedNumberOfVideoEncoderInstances
type GuaranteedNumberOfVideoEncoderInstances struct {
	TotalNumber int
	H264 int
}

// VideoSource
type VideoSource struct {
	Token string
	Framerate float64
	Resolution MediaBounds
	Imaging ImagingSettings
}

type ImagingSettings struct {
	BacklightCompensation BacklightCompensation
	Brightness float64
	ColorSaturation float64
	Contrast float64
	Exposure Exposure
	Focus FocusConfiguration
	IrCutFilter string //  'ON', 'OFF', 'AUTO'
	Sharpness float64
	WideDynamicRange WideDynamicRange
	WhiteBalance WhiteBalance
}

type BacklightCompensation struct {
	Mode string // 'ON' 'OFF'
	Level float64
}

type Exposure struct {
	Mode     string // 'AUTO', 'MANUAL'
	Priority string //  'LowNoise', 'FrameRate'
	Window   Rectangle
	MinExposureTime float64
	MaxExposureTime float64
	MinGain float64
	MaxGain float64
	MinIris float64
	MaxIris float64
	ExposureTime float64
	Gain float64
	Iris float64
}

type Rectangle struct {
	Top int
	Bottom int
	Left int
	Right int
}

type FocusConfiguration struct {
	AutoFocusMode string //  'AUTO', 'MANUAL'
	DefaultSpeed float64
	NearLimit float64
	FarLimit float64
}

type WideDynamicRange struct {
	Mode string //  'OFF', 'ON'
	Level float64
}

type WhiteBalance struct {
	Mode string //  'AUTO', 'MANUAL'
	CrGain float64
	CbGain float64
}

// Video Source Configuration
type VideoSourceConfiguration struct {
	Token string
	Name string
	SourceToken string
	Bounds IntRectangle
}

type IntRectangle struct {
	X int
	Y int
	Width int
	Height int
}

type IntRectangleRange struct {
	XRange IntRange
	YRange IntRange
	WidthRange IntRange
	HeightRange IntRange
}

type VideoSourceConfigurationOption struct {
	MaximumNumberOfProfiles int
	BoundsRange IntRectangleRange
	VideoSourceTokensAvailable string
}

type Multicast struct {
	Address IPAddress
	Port int
	TTL int
	AutoStart bool
}

type MetadataConfiguration struct {
	Token string
	Name string
	SessionTimeout string
	Multicast Multicast
}

type PTZStatusFilterOptions struct {
	PanTiltStatusSupported bool
	ZoomStatusSupported bool
	PanTiltPositionSupported bool
	ZoomPositionSupported bool
}

type MetadataConfigurationOptions struct {
	GeoLocation bool
	PTZStatusFilterOptions PTZStatusFilterOptions
}

type AudioSource struct {
	Token string
	Channels int //1: mono, 2: stereo
}

type AudioSourceConfiguration struct {
	Token string
	Name string
	SourceToken string
}

type AudioEncoderConfigurationOption struct {
	Encoding string // 'G711', 'G726', 'AAC'
	BitrateList int
	SampleRateList int
}