package database

import (
	internalMiddleware "github.com/nite-coder/blackbear-demo/internal/pkg/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "__gorm_span"
	callBackBeforeName = "telemetry:before"
	callBackAfterName  = "telemetry:after"
)

func before(db *gorm.DB) {
	tr := otel.Tracer("")
	ctx, span := tr.Start(db.Statement.Context, "gorm")

	lblRequestID := attribute.KeyValue{
		Key:   attribute.Key("request_id"),
		Value: attribute.StringValue(internalMiddleware.RequestIDFromContext(ctx)),
	}

	span.SetAttributes(lblRequestID)
	db.Statement.Context = ctx

	db.InstanceSet(gormSpanKey, span)
	return
}

func after(db *gorm.DB) {
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}

	span, ok := _span.(trace.Span)
	if !ok {
		return
	}
	defer span.End()

	// Error
	if db.Error != nil {
		lblError := attribute.KeyValue{
			Key:   attribute.Key("error"),
			Value: attribute.StringValue(db.Error.Error()),
		}

		evt := trace.WithAttributes(lblError)
		span.AddEvent("error", evt)

	}

	// sql
	lblSQL := attribute.KeyValue{
		Key:   attribute.Key("sql"),
		Value: attribute.StringValue(db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)),
	}
	span.SetAttributes(lblSQL)

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
