# TrackerAppBackend

TrackerAppBackend Codes
yep

## CORS Hatası ve Çözümü

### Karşılaşılan Hata

API'ye frontend tarafından yapılan isteklerde CORS (Cross-Origin Resource Sharing) hatası alındı. Terminalde aşağıdaki hata mesajları görüldü:

### Çözümün Açıklaması

- `Access-Control-Allow-Origin: "*"` - Tüm origin'lerden gelen isteklere izin verir. Üretim ortamında güvenlik için belirli domain'lere kısıtlanmalıdır.
- `Access-Control-Allow-Credentials: "true"` - Kimlik doğrulama bilgilerini içeren isteklere izin verir.
- `Access-Control-Allow-Headers` - İzin verilen HTTP başlıklarını belirtir.
- `Access-Control-Allow-Methods` - İzin verilen HTTP metodlarını belirtir.
- OPTIONS istekleri için özel işleme - Preflight isteklerini 204 (No Content) durum koduyla yanıtlar, bu da tarayıcıya asıl isteği göndermenin güvenli olduğunu bildirir.

Bu değişiklikler, frontend uygulamasının API'ye sorunsuz bir şekilde erişmesini sağlar.

### Güvenlik Notu

Üretim ortamında, `Access-Control-Allow-Origin` değerini `*` yerine frontend uygulamanızın gerçek domain'i ile değiştirmek daha güvenlidir:

```go
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://example.com"},
    AllowCredentials: true,
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    MaxAge: 24 * time.Hour,
}))
```

Bu yapılandırma, frontend uygulamanızın API'ye erişmesine izin verecektir.
