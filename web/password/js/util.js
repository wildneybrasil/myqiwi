util = {
    getLastSunday: function(date) {
        var crtDay = new Date(date).getUTCDay();
        var lstSunday = date - crtDay * 24 * 3600000;

        return lstSunday
    },

    getNextSaturday: function(date) {
        var crtDay = new Date(date).getUTCDay();
        var nxtSaturday = date + (6 - crtDay) * 24 * 3600000;

        return nxtSaturday
    },

    getPIndex: function() {
        return 0//mudar depois
    },

    readFilters: function(trigger) {
        var parent = $(trigger).parent();
        var vals = [];

        $(parent).find("input").each(function(index, elem) {
            vals.push($(elem).val());
        });

        $(parent).find("option:selected").each(function(index, elem) {
            vals.push($(elem).val());
        });

        console.log(vals);
        return vals
    },

    getMonthStr: function(targetName, option) {
        var targetVal = $("[name='" + targetName + "']").val();

        var date = new Date(targetVal);
        var year = date.getFullYear();
        var month = date.getMonth();
        var day = date.getDate();

        if(option == "start") {
            day = 1;
        } else if(option == "end") {
            day = new Date(year, month, 0).getDate();
        }

        var dateStr = year + "-" + (month + 2) + "-" + day;

        return dateStr
    },

    getSdate: function() {
        var targetVal = $("[name='from']").val();

        var date = new Date(targetVal);
        var year = date.getFullYear();
        var month = date.getMonth();
        var day = 1;

        var sdate = year + "-" + (month + 2) + "-" + day;

        return sdate
    },

    getEdate: function() {
        var targetVal = $("[name='until']").val();

        var date = new Date(targetVal);
        var year = date.getFullYear();
        var month = date.getMonth();
        var day = new Date(year, month, 0).getDate();

        var edate = year + "-" + (month + 2) + "-" + day;

        return edate
    },

    getCrntMonthStr: function(offset) {
        offset = typeof offset !== 'undefined' ? offset : 0;

        var CrntMonth = new Date();
        var year = CrntMonth.getFullYear();
        var month = CrntMonth.getMonth();

        var CrntMonthStr = year + "-" + (month + offset);

        return CrntMonthStr
    },

    getPrevMonth: function(date) {
        var newDate;
        var year = new Date(date).getFullYear();
        var month = new Date(date).getMonth();

        newDate = new Date(year, month - 1);

        return newDate
    },

    getNextMonth: function(date) {
        var newDate;
        var year = new Date(date).getFullYear();
        var month = new Date(date).getMonth();

        newDate = new Date(year, month + 2);

        return newDate
    },

    getPrevMonthStr: function(date) {
        var newDate = util.getPrevMonth(date);
        var newDateStr = format.dateToStr(newDate);

        return newDateStr
    },

    getNextMonthStr: function(date) {
        var newDate = util.getNextMonth(date);
        var newDateStr = format.dateToStr(newDate);

        return newDateStr
    },

    validateFilters: function(adjustFromDate) {
        adjustFromDate = typeof adjustFromDate !== 'undefined' ? adjustFromDate : true;

        var fromDate = format.strToDate($("[name='from']").val());
        var untilDate = format.strToDate($("[name='until']").val());

        if(fromDate >= untilDate) {
            if(adjustFromDate == true) {
                $("[name='from']").val(format.dateToStr(util.getPrevMonthStr(untilDate), "yyyy-mm"));
            } else {
                $("[name='until']").val(format.dateToStr(util.getNextMonthStr(fromDate), "yyyy-mm"));
            }
        }
    },

    getInput: function(input) {
        var _input = {}
        for(attr in input) {
            if(typeof input[attr] == "function") {
                _input[attr] = input[attr]();
            } else {
                _input[attr] = input[attr];
            }
        }
        return _input
    },

    resetForm: function(form, hide) {
        $(form).find("input").each(function(index, ele) {
            $(ele).val("");
        });
        if(hide == true) {
            $(form).hide();
            $(form).siblings("form").fadeIn();
        }
    }
},

format = {
    seriesData: function(data, id) {
        var parsedData = data.data;
        var seriesData;

        if(Object.prototype.toString.call(parsedData) === "[object Array]") {
            seriesData = [{
                name: "",
                id: id + "_" + 0,
                data: parsedData
            }];
            for(var i = 0; i < seriesData[0].data.length; i++) {
                seriesData[0].data[i][0] = format.strToDate(seriesData[0].data[i][0]);
                seriesData[0].data[i][1] = parseFloat(seriesData[0].data[i][1]);
            }
        } else {
            seriesData = [];
            var iterator = 0;
            for(var attr in parsedData) {
                seriesData.push({name: attr, id: attr + "_" + iterator, data: parsedData[attr]});
                iterator++;
            }
            for(var i = 0; i < seriesData.length; i++) {
                for(var j = 0; j < seriesData[i].data.length; j++) {
                    seriesData[i].data[j][0] = format.strToDate(seriesData[i].data[j][0]);
                    seriesData[i].data[j][1] = parseFloat(seriesData[i].data[j][1]);
                }
            }
        }

        return seriesData
    },

    tableData: function(data, formats, order, useAll) {
        formats = typeof formats !== 'undefined' ? formats : tables.defaultOpts.formats;
        order = typeof order !== 'undefined' ? order : tables.defaultOpts.order;
        useAll = typeof useAll !== 'undefined' ? useAll : tables.defaultOpts.useAll;

        var _data = data.data;
        var tableData = [];
        var row;


        for(var i in _data) {
            row = {};
            unordered = {};

            for (var j in order) {
                if(_data[i][order[j]] !== undefined) {
                    row[order[j]] = _data[i][order[j]];
                    delete _data[i][order[j]];
                }
            }

            if(useAll == true) {
                for (var j in _data[i]) {
                    row[j] = _data[i][j];
                }
            }

            tableData.push(row);
        }

        var j;

        for(var i = 0; i < tableData.length; i++) {
            for(var attr in tableData[i]) {
                if(formats[attr]) {
                    tableData[i][attr] = formats[attr](tableData[i][attr]);
                }
            }
        }
        
        return tableData
    },

    strToDate: function(str, inputFormat) {
        inputFormat = typeof inputFormat !== 'undefined' ? inputFormat : "dd/mm/yyyy";
        var strDate;
        var date;

        if(str && str !== null) {
            if(inputFormat == "dd/mm/yyyy") {
                strDate = str.split("/");
                if(strDate[0]) {
                    date = new Date(parseInt(strDate[2]), parseInt(strDate[1]), parseInt(strDate[0]));
                } else {
                    date = new Date(parseInt(strDate[2]), parseInt(strDate[1]));
                }
            }
            if(inputFormat == "yyyy-mm-dd") {
                strDate = str.split("-");
                if(strDate[2]) {
                    date = new Date(parseInt(strDate[0]), parseInt(strDate[1]), parseInt(strDate[2]));
                } else {
                    date = new Date(parseInt(strDate[0]), parseInt(strDate[1]));
                }
            }
        } else {
            return null
        }

        return date;
    },

    dateToStr: function(date, format) {
        format = typeof format !== 'undefined' ? format : "dd/mm/yyyy";
        var str;

        var year = new Date(date).getFullYear();
        var month = new Date(date).getMonth() + 1;
        var day = new Date(date).getDate();

        if(month <= 9) {
            month = "0" + month;
        }

        if(day <= 9) {
            day = "0" + day;
        }

        if(format == "dd/mm/yyyy") {
            str = day + "/" + month + "/" + year;
        }
        if(format == "yyyy-mm-dd") {
            str = year + "-" + month + "-" + day;
        }
        if(format == "yyyy-mm") {
            str = year + "-" + month;
        }

        return str
    },

    dateStr: function(dateStr, inputFormat, outputFormat) {
        inputFormat = typeof inputFormat !== 'undefined' ? inputFormat : "dd/mm/yyyy";
        outputFormat = typeof outputFormat !== 'undefined' ? outputFormat : "yyyy-mm-dd";

        var date = format.strToDate(dateStr, inputFormat);
        var str = format.dateToStr(date, outputFormat);

        return str
    }
},

terms = {
    nullVal: "Desconhecido",
    defaultCurrency: "R$"
}