package database

import (
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "__gorm_span"
	callBackBeforeName = "telemetry:before"
	callBackAfterName  = "telemetry:after"
)

func before(db *gorm.DB) {
	tr := global.Tracer("")
	ctx, span := tr.Start(db.Statement.Context, "gorm")
	span.SetAttribute("request_id", internalMiddleware.RequestIDFromContext(ctx))
	db.Statement.Context = ctx

	db.InstanceSet(gormSpanKey, span)
	return
}

func after(db *gorm.DB) {
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}
	ctx := db.Statement.Context

	span, ok := _span.(trace.Span)
	if !ok {
		return
	}
	defer span.End()

	// Error
	if db.Error != nil {
		span.AddEvent(ctx, "error", label.String("err", db.Error.Error()))
	}

	// sql
	span.SetAttribute("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))

	return
}

type TelemetryPlugin struct{}

func (tp *TelemetryPlugin) Name() string {
	return "telemetry_plugin"
}

func (tp *TelemetryPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TelemetryPlugin{}
