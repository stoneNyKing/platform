package router

import (
	"platform/hrsvc/imconf"
	"platform/hrsvc/services"
	"platform/mskit/grace"
	"platform/mskit/rest"
	"platform/mskit/trace"
	"platform/mskit/log"

	"platform/pfcomm/apis"
)

func InitRoute(prefix string,msapp *grace.MicroService) {

	// Logging domain.
	var logger =log.Mslog

	optracer,_:= apis.CreateTracer(imconf.Config.RecordAddr,imconf.Config.ServiceName,logger,imconf.Config.Debug,
		imconf.Config.ZipkinUrl,imconf.Config.AppdashAddr,imconf.Config.LightstepToken,imconf.Config.KafkaAddress)

	var options []trace.TraceOption
	options = append(options,trace.WithTracerOption(true))
	options = append(options,trace.OpenTracerOption(optracer))
	tracer := trace.NewTracer(options...)

	msapp.SetTracer(tracer)

	svc := services.StaffService{}
	deptsvc := services.DeptService{}
	staffsvc := services.DeptStaffService{}
	grpsvc := services.GroupService{}
	grpstaffsvc := services.GroupStaffService{}
	possvc := services.PositionService{}
	schedsvc := services.SchedulesService{}
	weeksvc := services.SchedWeeklyService{}
	monthsvc := services.SchedMonthlyService{}
	attsvc := services.AttendanceService{}
	applysvc := services.AttendanceApplyService{}
	plansvc := services.SchedPlanService{}
	org := services.OrganizationService{}

	mid := rest.RestMiddleware{Middle: LogMiddleware(logger), Object: logger}

	msapp.RegisterServiceWithTracer(prefix+"/org/:action", &org, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/org", &org, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/staff/:action", &svc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/staff", &svc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/dept/:action", &deptsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/dept", &deptsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/deptstaff/:action", &staffsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/deptstaff", &staffsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/group/:action", &grpsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/group", &grpsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/groupstaff/:action", &grpstaffsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/groupstaff", &grpstaffsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/position/:action", &possvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/position", &possvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/schedules/:action", &schedsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/schedules", &schedsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/weekly/:action", &weeksvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/weekly", &weeksvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/monthly/:action", &monthsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/monthly", &monthsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/attend/:action", &attsvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/attend", &attsvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/apply/:action", &applysvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/apply", &applysvc, tracer, logger, mid)

	msapp.RegisterServiceWithTracer(prefix+"/schedplan/:action", &plansvc, tracer, logger, mid)
	msapp.RegisterServiceWithTracer(prefix+"/schedplan", &plansvc, tracer, logger, mid)

	healthsvc := services.HealthCheckService{}

	hmid := rest.RestMiddleware{Middle: NoTokenCheck(logger), Object: logger}
	msapp.RegisterRestService(prefix+"/health", &healthsvc, hmid)
}
