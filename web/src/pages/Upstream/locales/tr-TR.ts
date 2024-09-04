/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
export default {
  'page.upstream.step.select.upstream': 'Upstream(Yukarı Akış)',
  'page.upstream.step.select.upstream.select.option': 'Gelişmiş',
  'page.upstream.step.select.upstream.select.option.serviceSelected':
    'Özel (Geçerli yapılandırma, bağlı hizmeti geçersiz kılar)',
  'page.upstream.step.select.upstream.select.none': 'Hiçbiri',
  'page.upstream.step.backend.server.domain.or.ip': 'Backend Sunucu Host/IP',
  'page.upstream.form.item-label.node.domain.or.ip': 'Hedef',
  'page.upstream.step.input.domain.name.or.ip': 'Domain veya IP girin',
  'page.upstream.step.valid.domain.name.or.ip': 'Geçerli bir domain veya IP adresi girin',
  'page.upstream.step.domain.name.or.ip': 'Hostname veya IP',
  'page.upstream.step.host': 'Host(istemci)',
  'page.upstream.step.input.port': 'Port numarası girin',
  'page.upstream.step.port': 'Port',
  'page.upstream.step.input.weight': 'Ağırlık girin',
  'page.upstream.step.weight': 'Ağırlık',
  'page.upstream.step.create': 'Oluştur',
  'page.upstream.step.name': 'Ad',
  'page.upstream.step.name.should.unique': 'Adınız benzersiz olmalıdır',
  'page.upstream.step.input.upstream.name': 'Upstream adı girin',
  'page.upstream.step.description': 'Açıklama',
  'page.upstream.step.input.description': 'Upstream için açıklama girin',
  'page.upstream.step.type': 'Algoritma',
  'page.upstream.step.pass-host': 'Host adı',
  'page.upstream.step.pass-host.pass': 'İstemci isteğinden gelen host adını bırak',
  'page.upstream.step.pass-host.node': 'Düğüm Listesinden etki alanını veya IPyi kullanın',
  'page.upstream.step.pass-host.rewrite': 'Gelişmiş Host (Gelecekte kullanımdan kaldırılacaktır)',
  'page.upstream.step.pass-host.upstream_host': 'Gelişmiş Host',
  'page.upstream.step.connect.timeout': 'Bağlantı zaman aşımı',
  'page.upstream.step.connect.timeout.desc':
    'İstekten yukarı akış sunucusuna bağlantı kurarken zaman aşımına uğradı.',
  'page.upstream.step.input.connect.timeout': 'Bağlantı zaman aşımı girin',
  'page.upstream.step.send.timeout': 'Gönderme zaman aşımı',
  'page.upstream.step.send.timeout.desc':
    'Yukarı akış sunucularına veri göndermek için zaman aşımı',
  'page.upstream.step.input.send.timeout': 'Zaman aşımı gönnderme girin',
  'page.upstream.step.read.timeout': 'Okuma zaman aşımı',
  'page.upstream.step.read.timeout.desc': 'Upstream sunucularının okuma zaman aşımı',
  'page.upstream.step.input.read.timeout': 'Okuma zaman aşımı girin',
  'page.upstream.step.healthyCheck.healthy.check': 'Sağlık kontrolü',
  'page.upstream.step.healthyCheck.healthy': 'Sağlıklı',
  'page.upstream.step.healthyCheck.unhealthy': 'Sağlıksız',
  'page.upstream.step.healthyCheck.healthy.status': 'Sağlık durumu',
  'page.upstream.step.healthyCheck.unhealthyStatus': 'Sağlıklı olmama durumu',
  'page.upstream.step.healthyCheck.active': 'Aktif',
  'page.upstream.step.healthyCheck.active.timeout': 'Zaman aşımı',
  'page.upstream.step.input.healthyCheck.active.timeout': 'Zaman aşımı girin',
  'page.upstream.step.input.healthyCheck.activeInterval': 'Aktif aralığı girin',
  'page.upstream.step.input.healthyCheck.active.req_headers': 'Request header girin',
  'page.upstream.step.healthyCheck.passive': 'Pasif',
  'page.upstream.step.healthyCheck.passive.http_statuses': 'HTTP durumları',
  'page.upstream.step.input.healthyCheck.passive.http_statuses': 'HTTP durumları girin',
  'page.upstream.step.healthyCheck.passive.tcp_failures': 'TCP hataları',
  'page.upstream.step.input.healthyCheck.passive.tcp_failures': 'TCP hataları girin',
  'page.upstream.step.keepalive_pool': 'Keepalive Pool',
  'page.upstream.notificationMessage.enableHealthCheckFirst':
    'Sağlık kontrolünü etkinleştirmek için önce bir upstream oluşturun',
  'page.upstream.upstream_host.required': 'Upstream host gerekli',

  'page.upstream.create': 'Yeni Upstream',
  'page.upstream.configure': 'Upstream Yapılandırma',
  'page.upstream.create.upstream.successfully': 'Upstream başarıyla oluşturuldu',
  'page.upstream.edit.upstream.successfully': 'Upstream başarıyla düzenlendi',
  'page.upstream.create.basic.info': 'Upstream Bilgileri',
  'page.upstream.create.preview': 'Önizleme',

  'page.upstream.list.id': 'ID',
  'page.upstream.list.name': 'Ad',
  'page.upstream.list.type': 'Tip',
  'page.upstream.list.description': 'Açıklama',
  'page.upstream.list.edit.time': 'Düzenleme zamanı',
  'page.upstream.list.operation': 'İşlemler',
  'page.upstream.list.edit': 'Düzenle',
  'page.upstream.list.confirm.delete': 'Upstream silmek istediğinize emin misiniz?',
  'page.upstream.list.confirm': 'Onay',
  'page.upstream.list.cancel': 'İptal',
  'page.upstream.list.delete.successfully': 'Upstream başarıyla silindi',
  'page.upstream.list.delete': 'Sil',
  'page.upstream.list': 'Upstream Listesi',
  'page.upstream.list.input': 'Lütfen girin',
  'page.upstream.list.create': 'Yarat',

  'page.upstream.type.roundrobin': 'Round Robin',
  'page.upstream.type.chash': 'CHash',
  'page.upstream.type.ewma': 'EWMA',
  'page.upstream.type.least_conn': 'Least conn',

  'page.upstream.list.content':
    'Yukarı akış listesi, oluşturulan yukarı akış hizmetlerini (yani arka uç hizmetleri) içerir ve yukarı akış hizmetlerinin birden çok hedef düğümünün yük dengelemesine ve sağlık denetimine izin verir.',

  'page.upstream.checks.active.timeout.description':
    'Aktif kontroller için soket zaman aşımı (saniye cinsinden)',
  'page.upstream.checks.active.unhealthy.interval.description':
    'Sağlıksız hedefler için kontroller arasındaki aralık (saniye cinsinden)',
  'page.upstream.checks.passive.healthy.http_statuses.description':
    'Hangi HTTP durumlarının başarılı olduğu düşünülmeli',
  'page.upstream.checks.passive.unhealthy.http_statuses.description':
    'Hangi HTTP durumları bir başarısızlık olarak kabul edilir?',
  'page.upstream.checks.passive.unhealthy.http_failures.description':
    'Bir hedefi sağlıksız olarak değerlendirmek için HTTP hatalarının sayısı',
  'page.upstream.checks.passive.unhealthy.tcp_failures.description':
    'Bir hedefi sağlıksız olarak değerlendirmek için TCP hatalarının sayısı',
  'page.upstream.scheme': 'Şema',

  'page.upstream.other.configuration.invalid': 'Lütfen upstream yapılandırmasını kontrol edin',

  'page.upstream.service-name.required': 'Gelişmiş server gerekli',
  'page.upstream.step.service-name': 'Server adı',
  'page.upstream.step.service-name.pass': 'İstemci isteğinden gelen server adını bırak',
  'page.upstream.step.service-name.rewrite': 'Gelişmiş Server',
  'page.upstream.step.service-name.custom_name': 'Gelişmiş Server',
};
