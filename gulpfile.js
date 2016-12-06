var gulp = require('gulp');
var browserify = require('browserify');
var babelify = require('babelify');
var source = require('vinyl-source-stream');
var browserSync = require('browser-sync').create();
var uglify = require('gulp-uglify');
var buffer = require('vinyl-buffer');
var minifyCss = require('gulp-minify-css');
var concatCss = require('gulp-concat-css');
var urlAdjuster = require('gulp-css-url-adjuster');

gulp.task('default', ['browserify','minify-css', 'copy-fonts'])

gulp.task('browserify', function () {
  return browserify('./source/app.js')
    .transform(babelify)
    .bundle()
    .pipe(source('bootreactor.js'))
    .pipe(buffer())
    .pipe(uglify())
    .pipe(gulp.dest('./public/'));
});

gulp.task('minify-css', function() {
  return gulp.src(['node_modules/bootstrap/dist/css/bootstrap.css'])
    .pipe(concatCss("min.css"))
    .pipe(urlAdjuster({
      replace:  ['../fonts/',''],
    }))
    .pipe(minifyCss({compatibility: 'ie8'}))
    .pipe(gulp.dest('./public/'));
});

gulp.task('copy-fonts', function() {
  return gulp.src('node_modules/bootstrap/dist/fonts/*')
    .pipe(gulp.dest('./public/'));
});

gulp.task('browser-sync', function() {
  browserSync.init({
    ui: {
      port: 8088
    },
    server: {
      baseDir: "public",
    }
  });
});
