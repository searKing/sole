{{- $my := .sole -}}
log:
  level: {{- default "info" $my.log_level | toJson -}}
  path: ./log/sole
  rotation_duration: 24h
  rotation_max_age: 168h
database:
  {{- with .databases.mysql.sole -}}
  #  dsn: "mysql://root:root@tcp(localhost:3306)/sole?max_conns=20&max_idle_conns=4"
  dsn: "mysql://{{ .user }}:{{ .password }}@tcp({{ .databases.mysql.sole.name }}|{{ .instances.mysql.single.host }}:{{ .instances.mysql.single.port }})/sole?max_conns=20&max_idle_conns=4"
  {{- else -}}
  dsn: memory
  {{- end -}}
  fail_after_duration: 5m
  max_wait_duration: 5s
system_secret: {{ default "Aa123456Aa123456" $my.system_secret | toJson }}
