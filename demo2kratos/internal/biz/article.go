package biz

import (
	"context"
	"log/slog"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/yylego/gormcnm"
	"github.com/yylego/gormrepo"
	"github.com/yylego/gormrepo/gormclass"
	"github.com/yylego/kratos-ebz/ebzkratos"
	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
	"github.com/yylego/kratos-examples/demo2kratos/internal/pkg/models"
	"github.com/yylego/kratos-gorm/gormkratos"
	"github.com/yylego/must"
	"gorm.io/gorm"
)

type Article struct {
	ID        int64
	Title     string
	Content   string
	StudentID int64
}

type ArticleUsecase struct {
	data *data.Data
	repo *gormrepo.Repo[models.T文章, *models.T文章Columns]
	log  *slog.Logger
}

func NewArticleUsecase(data *data.Data, logger *slog.Logger) *ArticleUsecase {
	return &ArticleUsecase{
		data: data,
		repo: gormrepo.NewRepo(gormclass.Use(&models.T文章{})),
		log:  logger,
	}
}

func (uc *ArticleUsecase) CreateArticle(ctx context.Context, a *Article) (*Article, *ebzkratos.Ebz) {
	must.Nice(a.Title)

	db := uc.data.DB()

	var v文章 *models.T文章

	if erk, err := gormkratos.Transaction(ctx, db, func(db *gorm.DB) *errors.Error {
		v文章 = &models.T文章{
			V标题:   a.Title,
			V内容:   a.Content,
			V学生ID: a.StudentID,
		}
		if err := uc.repo.With(ctx, db).Create(v文章); err != nil {
			return errors.New(500, "DB_ERROR", err.Error())
		}
		return nil
	}); err != nil {
		if erk != nil {
			return nil, ebzkratos.New(erk)
		}
		return nil, ebzkratos.New(pb.ErrorServerError("tx: %v", err))
	}
	return &Article{
		ID:        int64(v文章.ID),
		Title:     v文章.V标题,
		Content:   v文章.V内容,
		StudentID: v文章.V学生ID,
	}, nil
}

func (uc *ArticleUsecase) UpdateArticle(ctx context.Context, a *Article) (*Article, *ebzkratos.Ebz) {
	must.True(a.ID > 0)
	must.Nice(a.Title)

	db := uc.data.DB()

	if err := uc.repo.With(ctx, db).UpdatesM(func(db *gorm.DB, cls *models.T文章Columns) *gorm.DB {
		return db.Where(cls.ID.Eq(uint(a.ID)))
	}, func(cls *models.T文章Columns) gormcnm.ColumnValueMap {
		return cls.Kw(cls.V标题.Kv(a.Title)).Kw(cls.V内容.Kv(a.Content))
	}); err != nil {
		return nil, ebzkratos.New(pb.ErrorServerError("update: %v", err))
	}

	return a, nil
}

func (uc *ArticleUsecase) DeleteArticle(ctx context.Context, id int64) *ebzkratos.Ebz {
	must.True(id > 0)

	db := uc.data.DB()

	if err := uc.repo.With(ctx, db).DeleteW(func(db *gorm.DB, cls *models.T文章Columns) *gorm.DB {
		return db.Where(cls.ID.Eq(uint(id)))
	}); err != nil {
		return ebzkratos.New(pb.ErrorServerError("delete: %v", err))
	}
	return nil
}

func (uc *ArticleUsecase) GetArticle(ctx context.Context, id int64) (*Article, *ebzkratos.Ebz) {
	must.True(id > 0)

	db := uc.data.DB()

	v文章, erb := uc.repo.With(ctx, db).FirstE(func(db *gorm.DB, cls *models.T文章Columns) *gorm.DB {
		return db.Where(cls.ID.Eq(uint(id)))
	})
	if erb != nil {
		if erb.NotExist {
			return nil, ebzkratos.New(pb.ErrorServerError("not found: %v", erb.Cause))
		}
		return nil, ebzkratos.New(pb.ErrorServerError("db: %v", erb.Cause))
	}

	return &Article{
		ID:        int64(v文章.ID),
		Title:     v文章.V标题,
		Content:   v文章.V内容,
		StudentID: v文章.V学生ID,
	}, nil
}

func (uc *ArticleUsecase) ListArticles(ctx context.Context, page int32, pageSize int32) ([]*Article, int32, *ebzkratos.Ebz) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	db := uc.data.DB()

	// gormrepo FindPageAndCount replaces the stump's hand-written Count + Order + Offset + Limit
	// with one typed call that returns the current page plus the total row count together.
	// gormrepo 的 FindPageAndCount 把桩子里手写的 Count + Order + Offset + Limit
	// 收敛成一个类型安全的调用：一次拿到当页数据和总行数
	v文章们, total, err := uc.repo.With(ctx, db).FindPageAndCount(
		func(db *gorm.DB, cls *models.T文章Columns) *gorm.DB {
			return db
		},
		func(cls *models.T文章Columns) gormcnm.OrderByBottle {
			return cls.ID.Ob("asc")
		},
		&gormrepo.Pagination{
			Offset: int((page - 1) * pageSize),
			Limit:  int(pageSize),
		},
	)
	if err != nil {
		return nil, 0, ebzkratos.New(pb.ErrorServerError("list: %v", err))
	}

	items := make([]*Article, 0, len(v文章们))
	for _, v := range v文章们 {
		items = append(items, &Article{
			ID:        int64(v.ID),
			Title:     v.V标题,
			Content:   v.V内容,
			StudentID: v.V学生ID,
		})
	}
	return items, int32(total), nil
}

func (uc *ArticleUsecase) ListStudentArticles(ctx context.Context, studentID int64, page int32, pageSize int32) ([]*Article, int32, *ebzkratos.Ebz) {
	must.True(studentID > 0)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	db := uc.data.DB()

	// gormrepo FindPageAndCount 加类型安全的 student_id 过滤：带分页的关联查询
	// 仍是一个类型安全的调用，替掉桩子里手写的 Where + Count + Offset + Limit
	v文章们, total, err := uc.repo.With(ctx, db).FindPageAndCount(
		func(db *gorm.DB, cls *models.T文章Columns) *gorm.DB {
			return db.Where(cls.V学生ID.Eq(studentID))
		},
		func(cls *models.T文章Columns) gormcnm.OrderByBottle {
			return cls.ID.Ob("asc")
		},
		&gormrepo.Pagination{
			Offset: int((page - 1) * pageSize),
			Limit:  int(pageSize),
		},
	)
	if err != nil {
		return nil, 0, ebzkratos.New(pb.ErrorServerError("list student articles: %v", err))
	}

	items := make([]*Article, 0, len(v文章们))
	for _, v := range v文章们 {
		items = append(items, &Article{
			ID:        int64(v.ID),
			Title:     v.V标题,
			Content:   v.V内容,
			StudentID: v.V学生ID,
		})
	}
	return items, int32(total), nil
}
